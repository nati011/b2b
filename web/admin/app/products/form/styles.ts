// ProductsForm.styles.ts

import styled from "styled-components";
import { Button, Modal } from "antd";

export const StyledManageCategoriesButton = styled(Button)`
  margin: 8px;
  float: right;
  clear: both;
`;

export const StyledModal = styled(Modal)`
  .ant-modal-header {
    border-bottom: 1px solid #e8e8e8;
  }
  .ant-modal-title {
    font-weight: 600;
  }
  .ant-modal-body {
    padding: 24px;
  }
  .ant-modal-footer {
    border-top: 1px solid #e8e8e8;
  }
`;

export const StyledCategoryItem = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;

  &:last-child {
    border-bottom: none;
  }
`;

export const StyledAddCategoryButton = styled(Button)`
  float: right;
  margin-bottom: 16px;
`;
