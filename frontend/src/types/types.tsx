export interface Message {
  role: "user" | "assistant";
  content: string;
}

export interface StartInterviewResponse {
  session_id: string;
  message: string;
}

export interface ContinueInterviewRequest {
  message: string;
}

export interface ContinueInterviewResponse {
  message: string;
}
