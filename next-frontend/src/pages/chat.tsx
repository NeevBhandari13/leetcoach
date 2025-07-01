import React, { useEffect, useState } from 'react'; // react hook to manage state of components
import { useRouter } from 'next/router'; // next.js router hook to move between pages
import { Message, StartInterviewResponse } from '../types/types'; // import StartInterviewResponse type
import axios from '@/utils/axios';

const backendUrl = process.env.NEXT_PUBLIC_BACKEND_API_BASE_URL;

const ChatPage = () => {
    const router = useRouter(); // initialises router to be able to move between pages and access query parameters
    const { sessionID, reply } = router.query; // gets the session ID and reply from the query parameters

    // initialises messages state and setMessages function to manipulate messages
    // the state of a component always has a type (similar to an class' attribute)
    // the messages state is an array of Message objects
    const [messages, setMessages] = useState<Message[]>([])

    // useEffect to render an assistant message everytime there is a new response
    // will work when page is initialised from start Interview Response
    useEffect(() => {
        if (reply) {
            // new message object
            const initialMessage: Message = {
                role: 'assistant',
                content: String(reply)
            }
            setMessages([initialMessage]);
        }
    }, [reply])

    const handleSend = async (message: string) => {
        
    }

}

export default ChatPage;
