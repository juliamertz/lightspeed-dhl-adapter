use clap::{Parser, Subcommand};
use dhl::DHLError;

#[derive(Debug, Parser, Clone)]
pub struct Opts {
    #[clap(long, env = "DHL_USER_ID")]
    user_id: String,

    #[clap(long, env = "DHL_API_KEY")]
    api_key: String,

    #[clap(long, env = "DHL_ACCOUNT_ID")]
    account_id: String,

    #[command(subcommand)]
    command: Command,
}

#[derive(Subcommand, Clone, Debug)]
pub enum ShipmentCommand {
    Get {
        #[clap(short, long)]
        reference: String,
    },
}

#[derive(Subcommand, Clone, Debug)]
pub enum Command {
    Label {
        #[command(subcommand)]
        subcommand: ShipmentCommand,
    },
}

#[tokio::main]
pub async fn main() -> Result<(), DHLError> {
    let opts = Opts::parse();
    let credentials = dhl::Credentials {
        user_id: opts.user_id,
        api_key: opts.api_key,
        account_id: opts.account_id,
    };

    let client = dhl::DHLClient::new(credentials, false);
    client.authenticate().await?;

    match opts.command {
        Command::Label { subcommand } => match subcommand {
            ShipmentCommand::Get { reference } => {
                let label = client.get_label(&reference).await?;
                dbg!(label);
            }
        },
    };

    Ok(())
}
