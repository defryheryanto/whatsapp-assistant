# WhatsApp Assistant
Your WhatsApp Assistant

## Features
### WhatsApp Commands (Command Prefix: `%`)
1. `%commands`: Get All Commands
2. `%assign [role name] [@member1 @member2 @member3 ...]`: Assign role to mentioned members
3. `%[role name]`: Mention members of called role
4. `%all`: Mention all members in group
5. `%save [title] [content]`: Save the text
6. `%text [title]`: Get the saved text by title

## Development
### Setup
1. Install [Go](https://go.dev/doc/install)
2. Install [GCC](https://gcc.gnu.org/install/)
3. Install [Golang Migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#migrate-cli)
4. Run `go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

### GCC Install
1. Head to [WinLibs](https://winlibs.com/)
2. Go to download page and download the archive without LLVM
3. Extract and Add the bin folder to the environment variabbles
4. Verify installation by executing `gcc -v` in a new command promp tab

### Running the application
1. Clone the repository
2. Run `go run ./db/migrate.go`
3. Run `go run ./cmd/app/...` from the root project path
4. Copy the QRCode text from the terminal
5. Open https://www.the-qrcode-generator.com/ and navigate to 'Free Text' tab
6. Paste the QR Code text
7. Scan the generated QR from your WhatsApp
8. Done! Use commands or features from inside the chat

### How to use
1. Invite your WhatsApp Account (Scanned the generated QR) to any groups
2. Trigger the command by chat inside the group
3. Command will be triggered!

### Debugging the application (VS Code)
1. Add this code to your `.vscode/launch.json`
```
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/app"
        }
    ]
}
```
2. Press the debug button (F5 on windows)

### Adding new command
I'm making sure that it's easy to add a new command, you can refer to this steps:
1. Create a new go file under `./internal/assistant`
2. [Optional] Name it `command_[commandName].go`
3. In the file, define a struct that implements `assistant.commandAction` ([File here](https://github.com/defryheryanto/whatsapp-assistant/blob/6de82bb7b9befd4bcdf1475bf1e2543ee17a93ff/internal/assistant/commands.go#L26))
4. Define the execution process inside the `Execute` function
5. Add a new const defining the command ([File Here](https://github.com/defryheryanto/whatsapp-assistant/blob/6de82bb7b9befd4bcdf1475bf1e2543ee17a93ff/internal/assistant/commands.go#L16))
6. List the new command on `*assistant.WhatsappAssistant.getCommands()` ([File here](https://github.com/defryheryanto/whatsapp-assistant/blob/6de82bb7b9befd4bcdf1475bf1e2543ee17a93ff/internal/assistant/commands.go#L30))
7. Done! Your new command should be ready to use

## Disclaimer
It is important to note that utilizing a bot for WhatsApp can result in the suspension or banning of your WhatsApp account. Therefore, I strongly recommend refraining from using your primary WhatsApp account and using a secondary account instead to mitigate any potential risks<br>
Until now (19 June 2023) the possible reasons for your account getting banned are these:
1. `too many people blocked you`
2. `you sent too many messages to people who don't have you in their address books`
3. `you created too many groups with people who don't have you in their address books`
4. `you sent the same message to too many people`
5. `you sent too many messages to a broadcast list`
