use once_cell::sync::OnceCell;
use serde::{Deserialize, Serialize};

use std::path::Path;

#[derive(Clone, Serialize, Deserialize)]
pub struct Config {
    pub dhl: crate::dhl::models::Credentials,
    pub lightspeed: crate::lightspeed::models::Options,
    pub company_info: CompanyInfo,
}

impl std::fmt::Debug for Config {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        self.company_info.fmt(f)
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CompanyInfo {
    pub name: String,
    pub street: String,
    pub city: String,
    pub postal_code: String,
    pub country_code: String,
    pub number: String,
    pub addition: String,
    pub email: String,
    pub phone_number: String,
    pub personal_note: String,
}

#[derive(Debug, Default, Clone, Copy, Serialize, Deserialize, clap::ValueEnum)]
pub enum Environment {
    #[default]
    Development,
    Production,
}

static CONFIG: OnceCell<Config> = OnceCell::new();

pub fn load(path: impl AsRef<Path>) -> &'static Config {
    CONFIG.get_or_init(|| {
        let content = std::fs::read_to_string(path.as_ref()).unwrap();
        toml::from_str(&content).unwrap()
    })
}

pub fn get() -> &'static Config {
    CONFIG.get().expect("config to be initialized")
}
