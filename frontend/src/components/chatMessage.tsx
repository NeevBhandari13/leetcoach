import React from 'react';
import styles from "@/styles/ChatMessage.module.css";

// interface for chat message
interface ChatMessageProps {
    // role is either user or assistant
    role: "user" | "assistant";
    content: string;
}

// chat message component
const ChatMessage: React.FC<ChatMessageProps> = ({ role, content }: ChatMessageProps) => {

    const messageClass = `${styles.message} ${styles[role]}`;
    
    return (
        // the div tag has a style defined which defines how the chat bubbles will look and it contains the content which is what it will display
        <div
            className={messageClass}
        >
            {content} 
        </div>
    );
}

export default ChatMessage;
