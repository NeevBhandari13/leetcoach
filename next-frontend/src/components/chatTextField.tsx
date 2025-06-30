import React from "react";

interface ChatTextFieldProps {
    value: string;
    onChange: (event: React.ChangeEvent<HTMLInputElement>) => void; // allows the parent component to keep track of the input value as it is typed, every time there is a change to it
    disabled: boolean;
}

const ChatTextField: React.FC<ChatTextFieldProps> = ({ value, onChange, disabled }) => {
  return (
    // <input/> is a jsx text input field
    <input
      type="text" // type of input
      value={value} // sets the value of the input field to the value of the value prop
      onChange={onChange} // sets the onChange event handler to the onChange prop
      disabled={disabled} // sets the disabled attribute to the disabled prop
      placeholder="Type your message..."
      style={{ flex: 1, padding: '10px', fontSize: '16px' }}
    />
  );
};

export default ChatTextField;
