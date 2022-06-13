use clap::Parser;
use ssdownloader::{
    args::{Args, Commands},
    config::save_config,
};
fn main() {
    let args = Args::parse();

    // You can check for the existence of subcommands, and if found use their
    // matches just as you would the top level cmd
    match &args.command {
        Commands::Init {
            ss_api_key,
            ss_api_secret,
            zendesk_url,
        } => {
            save_config(ss_api_key, ss_api_secret, zendesk_url);
            println!("configuration saved")
        }
        Commands::Ticket { id, most_recent } => todo!(),
        Commands::Link { link } => todo!(),
    }
}
