use clap::{Parser, Subcommand};
use lightspeed::LightspeedError;

#[derive(Debug, Parser, Clone)]
pub struct Opts {
    #[clap(long, env = "LIGHTSPEED_KEY")]
    key: String,

    #[clap(long, env = "LIGHTSPEED_SECRET")]
    secret: String,

    #[clap(long, env = "LIGHTSPEED_CLUSTER")]
    cluster: String,

    #[command(subcommand)]
    command: Command,
}

#[derive(Subcommand, Clone, Debug)]
pub enum OrderCommand {
    Get {
        #[clap(short, long)]
        id: u64,
    },
}

#[derive(Subcommand, Clone, Debug)]
pub enum Command {
    Order {
        #[command(subcommand)]
        subcommand: OrderCommand,
    },
}

#[tokio::main]
pub async fn main() -> Result<(), LightspeedError> {
    let opts = Opts::parse();
    let credentials = lightspeed::Options {
        key: opts.key,
        secret: opts.secret,
        cluster: opts.cluster,
        // TODO: these options are only used in the adapter crate
        frontend: Default::default(),
        shop_id: Default::default(),
        cluster_id: Default::default(),
    };

    let client = lightspeed::LightspeedClient::new(credentials, false);

    match opts.command {
        Command::Order { subcommand } => match subcommand {
            OrderCommand::Get { id, } => {
                let order = client.get_order(id).await?;
                println!("{}", serde_json::to_string_pretty(&order)?);
            }
        },
    };

    Ok(())
}
