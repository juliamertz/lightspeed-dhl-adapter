use chrono::Local;
use reqwest::{Client as HttpClient, header::HeaderMap, StatusCode};
use thiserror::Error;
use tokio::sync::RwLock;
use tracing::info;

use crate::models::{ApiToken, Credentials, Draft, Label};

#[derive(Debug, Error)]
pub enum DHLError {
    #[error("invalid json: '{0}'")]
    Json(#[from] serde_json::Error),
    #[error("http request failed: '{0}'")]
    Http(#[from] reqwest::Error),
    #[error("invalid header value: '{0}'")]
    InvalidHeaderValue(#[from] reqwest::header::InvalidHeaderValue),
}

type Result<T> = std::result::Result<T, DHLError>;

const ENDPOINT: &str = "https://api-gw.dhlparcel.nl";

pub struct DHLClient {
    credentials: Credentials,
    token: RwLock<Option<ApiToken>>,
    http: HttpClient,
    skip_mutation: bool,
}

impl std::fmt::Debug for DHLClient {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.debug_struct("DHLClient")
            .field("account_id", &self.credentials.account_id)
            .field("skip_creation", &self.skip_mutation)
            .finish()
    }
}

impl DHLClient {
    pub fn new(opts: Credentials, skip_mutation: bool) -> Self {
        Self {
            credentials: opts,
            token: RwLock::new(None),
            http: reqwest::Client::new(),
            skip_mutation,
        }
    }

    async fn headers(&self) -> Result<HeaderMap> {
        let auth = self.authenticate().await?;
        let bearer = format!("Bearer {}", auth.access_token);

        let mut headers = HeaderMap::new();
        headers.insert("Authorization", bearer.try_into()?);
        headers.insert(
            "Content-Type",
            "application/json".try_into().expect("valid header"),
        );
        headers.insert(
            "Accept",
            "application/json".try_into().expect("valid header"),
        );
        Ok(headers)
    }

    pub async fn authenticate(&self) -> Result<ApiToken> {
        if let Some(token) = self.token.read().await.as_ref()
            && token.access_token_expiration - 15 > Local::now().timestamp()
        {
            return Ok(token.clone());
        }

        let mut token = self.token.write().await;

        let response: ApiToken = self
            .http
            .post(format!("{ENDPOINT}/authenticate/api-key"))
            .json(&self.credentials)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?;

        *token = Some(response.clone());

        Ok(response)
    }

    pub async fn get_drafts(&self) -> Result<Vec<Draft>> {
        let url = format!("{ENDPOINT}/drafts");
        let drafts = self
            .http
            .get(url)
            .headers(self.headers().await?)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?;

        Ok(drafts)
    }

    pub async fn create_draft(&self, draft: &Draft) -> Result<StatusCode> {
        if self.skip_mutation {
            info!(
                { draft = serde_json::to_string(&draft)? },
                "mocking draft creation"
            );
            return Ok(StatusCode::CREATED);
        }

        let url = format!("{ENDPOINT}/drafts");
        let response = self
            .http
            .post(url)
            .json(draft)
            .headers(self.headers().await?)
            .send()
            .await?
            .error_for_status()?;

        Ok(response.status())
    }

    pub async fn get_label(&self, reference: &str) -> Result<Option<Label>> {
        let url = format!("{ENDPOINT}/labels?orderReferenceFilter={reference}");
        let response: Vec<Label> = self
            .http
            .get(url)
            .headers(self.headers().await?)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?;

        let label = response.into_iter().next();
        Ok(label)
    }
}
