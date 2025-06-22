export interface Message {
  role: string;
  content: string;
}

export interface StartInterviewResponse {
  session_id: string;
  reply: string;
}

export interface ContinueInterviewRequest {
  sessionID: string;
  input: string;
}

export interface ContinueInterviewResponse {
  reply: string;
}
