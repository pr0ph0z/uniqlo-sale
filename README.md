## Uniqlo Sale  

A Golang-based scraper that tracks Uniqlo items on sale and sends the results to a Telegram channel or bot. It runs by comparing the previous fetched items with the current items and notifies you of any changes.

### Prerequisites  
- Go 1.21
- A Telegram bot token. You can get one from [BotFather](https://core.telegram.org/bots#botfather).  
- A [Telegram chat ID](https://gist.github.com/nafiesl/4ad622f344cd1dc3bb1ecbe468ff9f8a) where the results will be sent.  

### Setup  
1. Clone this repository:  
   ```bash  
   git clone https://github.com/pr0ph0z/uniqlo-sale.git  
   cd uniqlo-sale  
   ```  

2. Install dependencies:
   ```bash  
   go mod install
   ```  

3. Configure your environment variables:  
   Set the following environment variables in your shell or your shell script
   ```env  
   BOT_TOKEN=your_bot_token_here
   CHAT_ID=your_chat_id_here
   ```  

### Usage  

1. Run the scraper:  
   ```bash  
   go run main.go  
   ```  

2. The results will automatically be sent to your configured Telegram chat.  

### Limitation

Since each country has *almost* its own website, for the time being the scrapper only supports the Indonesian website (https://www.uniqlo.com/id/id/) and is hardcoded in project.

### Customization

To run the scraper periodically, you can set cron job on your own machine or you can use GitHub actions. I've already provided the workflow file in the `.github/workflows` directory so you can fork the repository directly, set the environment variables, and enable the workflow.

## Contributing  
Contributions are welcome! Feel free to submit a pull request or open an issue to suggest improvements.  
