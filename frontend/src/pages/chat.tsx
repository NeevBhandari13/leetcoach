import React, { useCallback, useEffect, useRef, useState } from 'react';
import { useRouter } from 'next/router';
import { Message } from '../types/types';
import ChatWindow from '@/components/chatWindow';
import ChatInput from '@/components/chatInput';
import Header from '@/components/Header';

const backendUrl = process.env.NEXT_PUBLIC_BACKEND_API_BASE_URL;

const ChatPage = () => {
    const router = useRouter();
    const { sessionID } = router.query;

    const [messages, setMessages] = useState<Message[]>([]);
    const [code, setCode] = useState('');
    const [leftWidth, setLeftWidth] = useState(50);
    const isDragging = useRef(false);
    const containerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (!sessionID) return;

        const initialText = sessionStorage.getItem('initialText');
        if (initialText) {
            setMessages([{ role: 'assistant', content: initialText }]);
            sessionStorage.removeItem('initialText');
            return;
        }

        fetch(`${backendUrl}/sessions/${sessionID}`)
            .then(res => res.json())
            .then(data => {
                const history: Message[] = (data.chat_history ?? []).map(
                    (m: { Role: string; Content: string }) => ({
                        role: m.Role as 'user' | 'assistant',
                        content: m.Content,
                    })
                );
                setMessages(history);
            })
            .catch(err => console.error('Error loading session:', err));
    }, [sessionID]);

    const handleSend = async (userMessage: string) => {
        setMessages(prev => [...prev, { role: 'user', content: userMessage }]);

        try {
            const response = await fetch(`${backendUrl}/sessions/${sessionID}/reply`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ message: userMessage, code }),
            });

            const data = await response.json();
            setMessages(prev => [...prev, { role: 'assistant', content: data.message }]);
        } catch (error) {
            console.error('Error sending message:', error);
        }
    };

    const handleMouseMove = useCallback((e: MouseEvent) => {
        if (!isDragging.current || !containerRef.current) return;
        const rect = containerRef.current.getBoundingClientRect();
        const newLeft = ((e.clientX - rect.left) / rect.width) * 100;
        setLeftWidth(Math.min(Math.max(newLeft, 20), 80));
    }, []);

    const handleMouseUp = useCallback(() => {
        isDragging.current = false;
    }, []);

    useEffect(() => {
        window.addEventListener('mousemove', handleMouseMove);
        window.addEventListener('mouseup', handleMouseUp);
        return () => {
            window.removeEventListener('mousemove', handleMouseMove);
            window.removeEventListener('mouseup', handleMouseUp);
        };
    }, [handleMouseMove, handleMouseUp]);

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
            <Header />
            <div
                ref={containerRef}
                style={{ display: 'flex', flex: 1, overflow: 'hidden', marginTop: '60px' }}
            >
                <div style={{ width: `${leftWidth}%`, display: 'flex', flexDirection: 'column', padding: '12px' }}>
                    <textarea
                        value={code}
                        onChange={e => setCode(e.target.value)}
                        placeholder="Write your code here..."
                        spellCheck={false}
                        style={{
                            flex: 1,
                            resize: 'none',
                            fontFamily: 'monospace',
                            fontSize: '14px',
                            padding: '12px',
                            border: '1px solid #ddd',
                            borderRadius: '8px',
                            outline: 'none',
                            backgroundColor: '#fafafa',
                        }}
                    />
                </div>

                <div
                    onMouseDown={() => { isDragging.current = true; }}
                    style={{
                        width: '5px',
                        cursor: 'col-resize',
                        backgroundColor: '#e0e0e0',
                        flexShrink: 0,
                        transition: 'background-color 0.15s',
                    }}
                    onMouseEnter={e => (e.currentTarget.style.backgroundColor = '#bdbdbd')}
                    onMouseLeave={e => (e.currentTarget.style.backgroundColor = '#e0e0e0')}
                />

                <div style={{ flex: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
                    <ChatWindow messages={messages} />
                    <ChatInput onSend={handleSend} />
                </div>
            </div>
        </div>
    );
};

export default ChatPage;
