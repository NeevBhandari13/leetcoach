import React from 'react';

// interface for chat message
interface ChatMessageProps {
    // role is either user or assistant
    role: "user" | "assistant";
    content: string;
}

// chat message component
const ChatMessage: React.FC<ChatMessageProps> = ({ role, content }: ChatMessageProps) => {

    const isUser = role === 'user';
    
    return (
        // the div tag has a style defined which defines how the chat bubbles will look and it contains the content which is what it will display
        <div
            style={{
                textAlign: isUser ? 'right' : 'left', // aligns content in chat bubble
                background: isUser ? '#f1f0f0' : "	#CBC3E3", // light grey for user, light purple for assistant
                margin: '8px 0', // having 2 values defaults to top-bottom and left-right, it creates space between chat bubbles
                padding: '10px 15px', // creates space inside content bubble for message content, 10px top and bottom, 15px left and right
                borderRadius: '10px', // rounds the corners
                maxWidth: '70%', // max width of a content bubble relative to chat window
                alignSelf: isUser ? 'flex-end' : 'flex-start' // aligns chat bubble to the right for user and left for assistant
            }}
        >
            {content} 
        </div>
    );
}

export default ChatMessage;
