import React, { useState } from "react";
import ChatTextField from './chatTextField';
import SendButton from './sendButton';

interface ChatInputProps {
    onSend: (input: string) => void;
}

const ChatInput: React.FC<ChatInputProps> = ({ onSend }) => {
    // creates the message state in component, initialises it to an empty string
    const [message, setMessage] = useState('');

    const handleSend = () => {
    if (message.trim() !== '') { // checks the message is not empty, trim gets rid of whitespace
      onSend(message); // calls the onSend function which is passed in as a prop
      setMessage(''); // sets the message back to ''
    }
    }

    return (
    <div style={{ display: 'flex', marginTop: '10px' }}>
        {/* Chat text field is the chat box, the message value is stored in this component */}
        {/* for onChange, e is the React.ChangeEvent<HTMLInputElement>, the target is the  native HTML <input> element rendered in the ChatTextField component and the value is the current input in it, setMessage is a function which updates the message state to the value in the input, we defined it above */}
      <ChatTextField value={message} onChange={(e) => setMessage(e.target.value)} disabled = {false} /> 
        {/* Passes down the handleSend function which is passed as a prop, disables the button if the message is empty */}
      <SendButton onClick={handleSend} disabled={!message.trim()} />
    </div>
  );
};


export default ChatInput;
