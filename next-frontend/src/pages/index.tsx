import StartInterviewButton from "@/components/startInterviewButton";

const backendUrl = process.env.NEXT_PUBLIC_BACKEND_API_BASE_URL;

// exports this as the main function, home means this will be at the root of the website
export default function Home() {

    // return is what should actually be rendered on the page
    return (
         <div style={{ textAlign: 'center', paddingTop: '100px' }}>
            <h1>LeetCoach</h1>
            <p>Welcome! Click the button below to start your interview.</p>
            <StartInterviewButton />
         </div>
    )

}
