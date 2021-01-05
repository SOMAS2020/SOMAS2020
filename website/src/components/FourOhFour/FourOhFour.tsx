import React from "react"
import { Button } from "react-bootstrap"
import { useHistory } from "react-router-dom"

const FourOhFour = () => {
  const history = useHistory()

  const goBack = () => {
    history.goBack()
  }

  const goHome = () => {
    history.push(`/`)
  }

  return (<div style={{ marginTop: 24 }}>
    <h1>Bear with Us</h1>
    <img src="https://mir-s3-cdn-cf.behance.net/project_modules/max_1200/26fa7351853877.58fc79a1747a2.gif" style={{ minWidth: `50vw` }} alt="Bear with us" />
    <p style={{ fontSize: `1.2em`, marginTop: 12 }}>We couldn't find what you're looking for; are you sure you have the right link?</p>
    <div>
      {
        history.length && <Button variant="dark" style={{ margin: 3 }} onClick={goBack}>Go Back</Button>
      }
      <Button variant="dark" style={{ margin: 3 }} onClick={goHome}>Go Home</Button>
    </div>
  </div>)
}

export default FourOhFour