import React, { useState } from "react";
import ChatTextField from './chatTextField';
import SendButton from './sendButton';
import styles from "@/styles/ChatInput.module.css";

interface ChatInputProps {
    onSend: (input: string) => void;
}

const ChatInput: React.FC<ChatInputProps> = ({ onSend }) => {
    // creates the message state in component, initialises it to an empty string
    const [message, setMessage] = useState('');

    // the target is the  native HTMLInputElement rendered in the ChatTextField component (which is specified) and the value is the current input in it, setMessage is a function which updates the message state to the value in the input, we defined it above
    const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        setMessage(e.target.value)
    }

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === 'Enter' && message.trim() !== '') {
        handleSend();
      }
    };

    /**
     * Checks if the message is not empty then calls the onSend function which is passed in as a prop and resets the message state to an empty string.
     * This function is triggered when the user clicks the send button.
     */
    const handleSend = () => {
    if (message.trim() !== '') { // checks the message is not empty, trim gets rid of whitespace
      onSend(message); // calls the onSend function which is passed in as a prop
      setMessage(''); // sets the message back to ''
    }
    }

    return (
    <div className={styles.chatInput}>
        {/* Chat text field is the chat box, the message value is stored in this component */}
        {/* onChange triggers the handleInput function everytime the input changes */}
      <ChatTextField value={message} onChange={handleInput} onKeyDown={handleKeyDown} disabled = {false} /> 
        {/* Passes down the handleSend function which is passed as a prop, disables the button if the message is empty */}
      <SendButton onClick={handleSend} disabled={!message.trim()} />
    </div>
  );
};


export default ChatInput;
