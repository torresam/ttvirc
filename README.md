# Twitch IRC

### **Setup**
1. Setup Environment Variables for connecting to IRC
    * Get an OAUTH token from [Twitch](https://twitchapps.com/tmi/)
    * Windows:
      * set TWITCH_USER=
      * set TWITCH_OAUTH=
    * Linux:
      * export TWITCH_USER=
      * export TWITCH_OAUTH=

---


Connects to Twitch IRC using SSL endpoint: `irc.chat.twitch.tv:6697`  
Listens on `localhost:8000` for command requests 
  * Get a list of channels currently connected to:  
    ```
    curl localhost:8000/channels
    ``` 
  * Connect to the specified channel:  
    ```
    curl -x PUT localhost:8000/join/{channel}
    ``` 
  * Leave the specified channel:
    ```
    curl -x PUT localhost:8000/leave/{channel}
    ```