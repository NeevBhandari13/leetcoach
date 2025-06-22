import React, { useState } from 'react'; // react hook to manage state of components
import { useRouter } from 'next/router'; // next.js router hook to move between pages
import { StartInterviewResponse } from '../types/types'; // import StartInterviewResponse type


const baseUrl = process.env.BACKEND_API_BASE_URL;

// exports this as the main function, home means this will be at the root of the website
export default function Home() {
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
            const startInterviewURL = `${baseUrl}/start-interview`; // url to start interview
            // call start interview api
            const response = await fetch(startInterviewURL, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });
            
            // handle error by checking response
            if (!response.ok) {
                throw new Error('Failed to start interview');
            }
            
            // response.json() returns a promise that resolves to the JSON data
            // the await keyword is used to wait for the promise to resolve so we get the data rather than the promise
            // the data is then assigned to the StartInterviewResponse type in the data variable
            const data: StartInterviewResponse = await response.json();

            // router.push allows us to move to the chat page with the sessionID and response as query parameters
            // the query parameters means we store them and have them available when the page loads
            router.push({
                pathname: '/chat',
                query: {
                sessionID: data.session_id,
                response: data.reply,
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

}
