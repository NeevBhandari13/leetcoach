import React, { useState } from 'react'; // react hook to manage state of components
import { useRouter } from 'next/router'; // next.js router hook to move between pages
import { StartInterviewResponse } from '@/types/types'; // import StartInterviewResponse type

const backendUrl = process.env.NEXT_PUBLIC_BACKEND_API_BASE_URL;

// exports this as the main function, home means this will be at the root of the website
export default function StartInterviewButton() {
    // initialises loading state and setLoading function to manipulate loading state
    const [loading, setLoading] = useState(false);
    // initialises router to be ablet to move between pages
    const router = useRouter();

    // triggered when start interview button is pressed
    // function to handle the start of the interview
    // using async means we can use await in it for things like calling the backend
    async function handleStartInterview() {
        // sets loading state to true because button has been pressed
        setLoading(true);

        // try lets us handle errors gracefully if something goes wrong with the network calls etc.
        try {
            const startInterviewURL = `${backendUrl}/start-interview`; // url to start interview
            // call start interview api

            const response = await fetch(startInterviewURL, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            })
            
            // response.json() returns a promise that resolves to the JSON data
            // the await keyword is used to wait for the promise to resolve so we get the data rather than the promise
            // the data is then assigned to the StartInterviewResponse type in the data variable
            const data = await response.json();
            
            
            // router.push allows us to move to the chat page with the sessionID and response as query parameters
            // the query parameters means we store them and have them available when the page loads
            router.push({
                pathname: '/chat',
                query: {
                sessionID: data.sessionId,
                initialText: data.reply,
                },
            });

        } catch (error) {
            // catches error we throw if interview doesn't start
            console.error('Error starting interview:', error);
        } finally {
            // sets loading state to false
            // important so the button becomes clickable again
            setLoading(false);
        }
        
        
    }

    // return is what should actually be rendered on the page
    return (
            <button onClick={handleStartInterview} disabled={loading}>
                {/* button that calls handleStartInterview when clicked */}
                {/* disabled = {loading} means the button is disabled while loading */}
                {/* if loading is true show this text*/}
                {loading ? 'Starting...' : 'Start Interview'}
            </button>
    )

}
