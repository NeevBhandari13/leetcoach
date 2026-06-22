import React, { useState } from 'react';
import { useRouter } from 'next/router';
import { StartInterviewResponse } from '@/types/types';
import styles from '@/styles/startInterviewButton.module.css';

const backendUrl = process.env.NEXT_PUBLIC_BACKEND_API_BASE_URL;

export default function StartInterviewButton() {
    const [loading, setLoading] = useState(false);
    const router = useRouter();

    async function handleStartInterview() {
        setLoading(true);
        try {
            const response = await fetch(`${backendUrl}/start`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
            });

            const data: StartInterviewResponse = await response.json();

            sessionStorage.setItem('initialText', data.message);
            router.push({
                pathname: '/chat',
                query: { sessionID: data.session_id },
            });
        } catch (error) {
            console.error('Error starting interview:', error);
        } finally {
            setLoading(false);
        }
    }

    return (
        <button className={styles.button} onClick={handleStartInterview} disabled={loading}>
            {loading ? 'Starting...' : 'Start Interview'}
        </button>
    );
}
