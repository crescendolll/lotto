import React, { useEffect, useState } from 'react'
import { Header } from '../../layout/Header/Header'
import { Footer } from '../../layout/Footer/Footer'
import { TippsTabelle } from '../../forms/TippsTabelle/TippsTabelle'
import { PlusButton } from '../../forms/PlusButton/PlusButton'
import { AnsichtSwitch } from '../../forms/AnsichtSwitch/AnsichtSwitch'
import { getTips, getOpenGames, getClosedGames, submitTip } from '../../../api'
import { Button } from '../../Button'
import DayPickerInput from 'react-day-picker/DayPickerInput'
import 'react-day-picker/lib/style.css'

export const TippsPage = ({ auth, onLogout }) => {
  const [tiplist, setTips] = useState([])
  const [openGames, setOpenGames] = useState([])
  const [date, setDate] = useState('')
  const [tipsubmission, setTipsubmission] = useState('')
  const [showTipSubmit, setShowTipSubmit] = useState(false)

  useEffect(() => {
    getTips(auth, '', '').then((data) => {
      return setTips(data.statistik)
    })
  }, [])

  const onGetOpenGames = (auth) => {
    console.log(auth)
    getOpenGames(auth).then((data) => {
      return setOpenGames(data.ziehungstage)
    })
  }

  const onOpenSubmitTip = () => {
    setShowTipSubmit(true)
  }

  const onSubmit = (e) => {
    e.preventDefault()
    if (date && tipsubmission) {
      submitTip(auth, tipsubmission, date.toLocaleDateString('en-CA')).then(
        (data) => console.log(data)
      )
    }
    setShowTipSubmit(false)
    setDate('')
    setTipsubmission('')
  }

  const onChangeDate = (date) => {
    setDate(date)
  }

  const onChangeTipsubmission = (event) => {
    setTipsubmission(event.target.value)
  }

  return (
    <div>
      <Header auth={auth} onLogout={onLogout} />
      <Button onClick={() => onGetOpenGames(auth)}>get open Games</Button>
      {openGames.length > 0 ? <TippsTabelle tips={openGames} /> : null}
      {showTipSubmit ? (
        <form
          onSubmit={onSubmit}
          className='my-9 bg-grau max-w-min min-w-max p-2'
        >
          <Button>Submit</Button>
          <DayPickerInput
            onDayChange={(date) => onChangeDate(date)}
            placeholder='Available Dates'
            hideOnDayClick='true'
            dayPickerProps={{
              disabledDays: {
                before: new Date(),
              },
            }}
          />
          <input
            type='text'
            size='20'
            placeholder='Enter Tipp'
            onChange={onChangeTipsubmission}
            className=' mx-2 border rounded-md border-dunkelgrau'
          />
        </form>
      ) : (
        <div className='bg-grau p-2 max-w-min min-w-max my-9'>

          <PlusButton onClick={onOpenSubmitTip}/>Submit a Tip
        </div>
      )}
      {/* <AnsichtSwitch /> */}
      <TippsTabelle tips={tiplist} />
      <Footer />
    </div>
  )
}
