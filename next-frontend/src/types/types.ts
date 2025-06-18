export interface Message {
    role: string;
    content: string;
}

export interface StartInterviewResponse {
    session_id: string;
    reply: Message[];
}

export interface ContinueInterviewRequest {
    session_id: string;
    input: string;
}

export interface ContinueInterviewResponse {
    reply: string;
}
