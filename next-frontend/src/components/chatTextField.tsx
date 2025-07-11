import React from "react";
import styles from "@/styles/chatTextField.module.css";

interface ChatTextFieldProps {
    value: string;
    onChange: (event: React.ChangeEvent<HTMLInputElement>) => void; // allows the parent component to keep track of the input value as it is typed, every time there is a change to it
    onKeyDown?: (e: React.KeyboardEvent<HTMLInputElement>) => void; // handled by this component because it is listening for keyboard events
    disabled: boolean;
}

const ChatTextField: React.FC<ChatTextFieldProps> = ({ value, onChange, onKeyDown, disabled }) => {
  return (
    // <input/> is a jsx text input field
    <input
      type="text" // type of input
      value={value} // sets the value of the input field to the value of the value prop
      className={styles.inputField}
      onChange={onChange} // sets the onChange event handler to the onChange prop
      onKeyDown={onKeyDown} // tracks if the enter key is pressed
      disabled={disabled} // sets the disabled attribute to the disabled prop
      placeholder="Type your message..."
    />
  );
};

export default ChatTextField;
