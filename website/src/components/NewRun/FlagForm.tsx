import React from "react"
import { Flag } from "../../wasmAPI"
import { Col, OverlayTrigger, Tooltip, Form } from 'react-bootstrap'

type Props = {
  flag: Flag,
  setFlag: (val: string) => Promise<void>,
  disabled: boolean,
}

const FlagForm = (props: Props) => {
  const { flag, setFlag, disabled } = props

  const handleChange = async (event: React.ChangeEvent<any>) => {
    await setFlag(event.target.value)
  }
  return <Col xs={4}>
    <Form>
      <Form.Group>
        <Form.Label>
          <OverlayTrigger
            placement="top"
            overlay={
              <Tooltip id={flag.Name}>
                {flag.Usage} (Type: {flag.Type})
              </Tooltip>
            }
          >
            <span style={{ wordBreak: `break-all` }}>{flag.Name}</span>
          </OverlayTrigger >
        </Form.Label>
        <Form.Control value={flag.Value} onChange={handleChange} readOnly={disabled} isInvalid={flag.InvalidReason !== undefined} />
      </Form.Group>
    </Form>
  </Col >
}

export default FlagForm