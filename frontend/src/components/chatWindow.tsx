import React from "react";
import ChatMessage from "./chatMessage";
import { Message } from "../types/types";
import styles from "@/styles/ChatWindow.module.css";

// interface for chat window
interface ChatWindowProps {
    // array of messages
    messages: Message[];
}


const ChatWindow: React.FC<ChatWindowProps> = ({ messages }) => {
    return (
        // 70vh, means max height is 70% of the viewport
        // overflowY: 'scroll'
        // padding: '10px' means 10px on the top and bottom
        <div className={styles.chatWindow}>
            {/* map through the array of messages and render a chat message for each one */}
            {/* index is automatically provided by the map function, the index of the message in the array, allows us to avoid rerendering the same message every time since we are only appending so messages will have stable keys which can be reused */}
            {messages.map((message, index) => (
                <ChatMessage key={index} role={message.role} content={message.content} />
            ))}
        </div>
    )
    
}

export default ChatWindow;

