export interface Message {
  role: "user" | "assistant";
  content: string;
}

export interface StartInterviewResponse {
  sessionId: string;
  initialText: string;
}

export interface ContinueInterviewRequest {
  sessionId: string;
  input: string;
}

export interface ContinueInterviewResponse {
  reply: string;
}
