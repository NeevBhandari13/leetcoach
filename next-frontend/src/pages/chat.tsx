import React, { useEffect, useState } from 'react'; // react hook to manage state of components
import { useRouter } from 'next/router'; // next.js router hook to move between pages
import { ContinueInterviewRequest, Message, StartInterviewResponse } from '../types/types'; // import StartInterviewResponse type
import axios from '@/utils/axios';
import ChatWindow from '@/components/chatWindow';
import SendButton from '@/components/sendButton';

const backendUrl = process.env.NEXT_PUBLIC_BACKEND_API_BASE_URL;

const ChatPage = () => {
    const router = useRouter(); // initialises router to be able to move between pages and access query parameters
    const { sessionID, responseText } = router.query; // gets the session ID and responseText from the query parameters

    // initialises messages state and setMessages function to manipulate messages
    // the state of a component always has a type (similar to an class' attribute)
    // the messages state is an array of Message objects
    const [messages, setMessages] = useState<Message[]>([])

    // useEffect to render an assistant message everytime there is a new response
    // will work when page is initialised from start Interview Response
    useEffect(() => {
        if (responseText) {
            // new message object
            const initialMessage: Message = {
                role: 'assistant',
                content: String(responseText)
            }
            setMessages([initialMessage]);
        }
    }, [responseText])

    const handleSend = async (message: string) => {
        // add new message to our messages state
        setMessages(prevMessages => [...prevMessages, { role: 'user', content: message }]);

        try {
            // package message into continue interview request
            const continueInterviewRequest: ContinueInterviewRequest = {
                sessionID: String(sessionID),
                input: message
            }
            // send message to backend
            // directly send ContinueInterviewRequest
            const response = await fetch(`${backendUrl}/continue-interview`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(continueInterviewRequest),
            })
            
            // response.json() parses the response body as JSON into a JavaScript object
            const data = await response.json();

            setMessages(prevMessages => [...prevMessages, { role: 'assistant', content: data.response_text}]);

    } catch (error) {
        console.error("Error sending message:", error);
    }

    return (
        <div style={{ textAlign: 'center', paddingTop: '100px' }}>
            <h1>LeetCoach</h1>
            <div>
                
            </div>
            <ChatWindow />
            <SendButton onClick={() => handleSend('hello')} disabled={false} />
        </div>
    )
}

}

export default ChatPage;
