import React from "react";
import styles from "@/styles/sendButton.module.css";

interface SendButtonProps {
    onClick: () => void;
    disabled: boolean;
}

const SendButton: React.FC<SendButtonProps> = ({ onClick, disabled }) => {
    return (
        <button
            className={styles.button}
            onClick={onClick}
            disabled={disabled}
        >
            Send
        </button>
    );
};

export default SendButton;
