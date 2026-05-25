import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import { Message } from '../types/types';
import ChatWindow from '@/components/chatWindow';
import ChatInput from '@/components/chatInput';
import Header from '@/components/Header';

const backendUrl = process.env.NEXT_PUBLIC_BACKEND_API_BASE_URL;

const ChatPage = () => {
    const router = useRouter();
    const { sessionID, initialText } = router.query;

    const [messages, setMessages] = useState<Message[]>([]);

    useEffect(() => {
        if (initialText) {
            setMessages([{ role: 'assistant', content: String(initialText) }]);
        }
    }, [initialText]);

    const handleSend = async (userMessage: string) => {
        setMessages(prev => [...prev, { role: 'user', content: userMessage }]);

        try {
            const response = await fetch(`${backendUrl}/sessions/${sessionID}/reply`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ message: userMessage }),
            });

            const data = await response.json();

            setMessages(prev => [...prev, { role: 'assistant', content: data.message }]);
        } catch (error) {
            console.error('Error sending message:', error);
        }
    };

    return (
        <div style={{ textAlign: 'center', paddingTop: '100px' }}>
            <Header />
            <ChatWindow messages={messages} />
            <ChatInput onSend={handleSend} />
        </div>
    );
};

export default ChatPage;
