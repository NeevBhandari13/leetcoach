import { Send } from "lucide";
import React from "react";

interface SendButtonProps {
    onClick: () => void;
    disabled: boolean;
}

const SendButton: React.FC<SendButtonProps> = ({ onClick, disabled }) => {

    return (
        <button 
            onClick={onClick} 
            disabled={disabled}
            style={{ padding: '10px 20px', marginLeft: '10px', fontSize: '16px' }}
        >
            Send
        </button>
    )
}

export default SendButton;
