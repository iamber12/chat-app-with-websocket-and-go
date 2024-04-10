document.addEventListener('DOMContentLoaded', function() {
    var ws; // WebSocket connection
    var email; // User email
    var loginForm = document.getElementById('loginForm');
    var emailInput = document.getElementById('emailInput');
    var loginButton = document.getElementById('loginButton');
    var chatContainer = document.getElementById('chatContainer');
    var messageBox = document.getElementById('messageBox');
    var messageInput = document.getElementById('messageInput');
    var sendButton = document.getElementById('sendButton');

    function showMessage(message) {
        var messageElement = document.createElement('div');
        messageElement.textContent = message.From + ': ' + message.Message;
        messageBox.appendChild(messageElement);
        messageBox.scrollTop = messageBox.scrollHeight; // Auto-scroll to the latest message
    }

    loginButton.addEventListener('click', function() {
        email = emailInput.value.trim();
        if (email) {
            // Hide the login form and show the chat container
            loginForm.style.display = 'none';
            chatContainer.style.display = 'block';
            connect(); // Establish WebSocket connection after logging in
        } else {
            alert('Please enter your email.');
        }
    });

    function connect() {
        ws = new WebSocket('ws://localhost:8000/ws');

        ws.onopen = function() {
            console.log('Connected to the chat server');
        };

        ws.onmessage = function(event) {
            var message = JSON.parse(event.data);
            showMessage(message);
        };

        ws.onclose = function(e) {
            console.log('WebSocket is closed. Reconnect will be attempted in 1 second.', e.reason);
            setTimeout(function() {
                connect();
            }, 1000);
        };

        ws.onerror = function(err) {
            console.error('WebSocket encountered error: ', err.message, 'Closing WebSocket');
            ws.close();
        };
    }

    sendButton.addEventListener('click', function() {
        if (!ws || ws.readyState !== WebSocket.OPEN) {
            console.error("WebSocket isn't open. Unable to send message.");
            return;
        }
        var message = {
            from: email, // User's email as the message sender
            message: messageInput.value
        };
        ws.send(JSON.stringify(message));
        messageInput.value = ''; // Clear the input after sending
    });
});
