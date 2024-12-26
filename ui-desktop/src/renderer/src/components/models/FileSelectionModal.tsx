//import Modal from '../../contracts/modals/Modal';
import Modal from '../contracts/modals/Modal';
import styled from 'styled-components';
import {
    TitleWrapper,
    Title,
    RightBtn
} from '../contracts/modals/CreateContractModal.styles';
import { useState } from 'react';
import Form from 'react-bootstrap/Form';
import { Sp } from '../common'

const bodyProps = {
    height: '750px',
    width: '70%',
    maxWidth: '100%',
    overflow: 'hidden',
    onClick: e => e.stopPropagation()
}

const RowContainer = styled.div`
  padding: 1rem;
    border: 1px solid;
    background-color: rgba(0,0,0,0.4);
    margin-bottom: 1rem;
  `

const FileSelectionModal = ({ isActive, handleClose }) => {

    if (!isActive) {
        return <></>;
    }

    const [files, setFiles] = useState<any>([]);

    return (
        <Modal
            onClose={() => {
                handleClose();
            }}
            bodyProps={bodyProps}
        >
            <TitleWrapper>
                <Title>Select Files</Title>
            </TitleWrapper>

            <Sp mt={2}>
                <Form.Group controlId="formFile" className="mb-3">
                    <Form.Label>Set desired model name</Form.Label>
                    <Form.Control type="text" />
                </Form.Group>
            </Sp>
            <Sp mt={2}>
                <Form.Group controlId="formFile" className="mb-3">
                    <Form.Label>Select files required to run model (including .gguf)</Form.Label>
                    <Form.Control type="file" multiple onChange={(e => {
                        setFiles(Object.values((e.currentTarget as any).files))
                    })} />
                </Form.Group>
            </Sp>

            {
                !files.length ? null :
                    (
                        files.map(f => {
                            return <RowContainer>
                                <div>Name: {f.name}</div>
                                <div>Path: {f.path}</div>
                                <div>Size: {(f.size / 1024).toFixed(0)}</div>
                            </RowContainer>
                        })
                    )
            }

            <Sp mt={2} style={{ dispay: 'flex', justifyContent: 'center'}}>
                <RightBtn>Pin Model Files</RightBtn>
            </Sp>

        </Modal>
    );
}

export default FileSelectionModal;