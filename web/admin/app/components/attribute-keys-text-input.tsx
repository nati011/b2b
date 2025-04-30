import { useState } from 'react';
import { Form, Input, Button, Tag } from 'antd';

const AttributeKeysInput = () => {
    const [attributeKeys, setAttributeKeys] = useState<string[]>([]);
    const [inputValue, setInputValue] = useState('');

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setInputValue(e.target.value);
    };

    const handleInputKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter' || e.key === ',') {
            e.preventDefault(); // Prevent form submission on Enter
            const newKey = inputValue.trim();
            if (newKey && !attributeKeys.includes(newKey)) {
                setAttributeKeys([...attributeKeys, newKey]);
            }
            setInputValue('');
        }
    };

    const handleRemoveKey = (keyToRemove: string) => {
        setAttributeKeys(attributeKeys.filter((key) => key !== keyToRemove));
    };

    return (
        <div>
            <Input
                name="attributeKeys" // Add the name attribute
                placeholder="Enter attribute keys (comma or Enter separated)"
                value={inputValue}
                onChange={handleInputChange}
                onKeyDown={handleInputKeyDown}
            />
            <div style={{ marginTop: '8px' }}>
                {attributeKeys.map((key) => (
                    <Tag key={key} closable onClose={() => handleRemoveKey(key)}>
                        {key}
                    </Tag>
                ))}
            </div>
        </div>
    );
};

export default AttributeKeysInput;