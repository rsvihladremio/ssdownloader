use clap::{Parser, Subcommand};

/// Args for ssdownloader
#[derive(Parser, Debug)]
#[clap(
    author = "Ryan Svihla",
    version = "0.1.0",
    about = "sendsafely file downloader to automate downloading all files from a zendesk ticket",
    long_about = "it is often painful to download the necessary data to analyze a long running ticket in zendesk, so ssdownloader will retrieve all of the files when calling a particular ticket number"
)]
#[clap(propagate_version = true)]
pub struct Args {
    #[clap(subcommand)]
    pub command: Commands,
}
#[derive(Subcommand, Debug)]
pub enum Commands {
    /// Initializes configuration
    Init {
        /// sendsafely api key
        #[clap(help = "sendsafely key to use")]
        ss_api_key: String,
        /// sendsafely api secret
        #[clap(help = "sendsafely api secret to use")]
        ss_api_secret: String,
        #[clap(help = "zendesk url to use for ticket downloads")]
        zendesk_url: String,
    },
    Ticket {
        #[clap(help = "ticket number for the sendsafely link")]
        id: Option<i64>,
        /// most recent posts limited to x most recent posts
        #[clap(short, long)]
        most_recent: Option<u8>,
    },
    Link {
        #[clap(help = "url for the sendsafely link will download all files in that packages")]
        link: String,
    },
}
