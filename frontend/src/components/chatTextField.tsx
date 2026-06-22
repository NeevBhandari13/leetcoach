import React, { useEffect, useRef } from "react";
import styles from "@/styles/chatTextField.module.css";

interface ChatTextFieldProps {
    value: string;
    onChange: (event: React.ChangeEvent<HTMLTextAreaElement>) => void;
    onKeyDown?: (e: React.KeyboardEvent<HTMLTextAreaElement>) => void;
    disabled: boolean;
}

const ChatTextField: React.FC<ChatTextFieldProps> = ({ value, onChange, onKeyDown, disabled }) => {
  const ref = useRef<HTMLTextAreaElement>(null);

  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    onChange(e);
    if (ref.current) {
      ref.current.style.height = 'auto';
      ref.current.style.height = `${Math.min(ref.current.scrollHeight, 160)}px`;
    }
  };

  useEffect(() => {
    if (value === '' && ref.current) {
      ref.current.style.height = 'auto';
    }
  }, [value]);

  return (
    <textarea
      ref={ref}
      value={value}
      className={styles.inputField}
      onChange={handleChange}
      onKeyDown={onKeyDown}
      disabled={disabled}
      placeholder="Type your message..."
      rows={1}
    />
  );
};

export default ChatTextField;
