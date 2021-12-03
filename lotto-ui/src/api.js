const request = (data) => {
  return fetch('/', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  }).then((response) => response.json())
}

//Sign up player with name, password (NO AUTH) --> RESPONSE: auth: ""
export const signUp = (name, password) => {
  return request({
    methode: 'registriere',
    param: {
      name,
      passwort: password,
    },
  })
}

//Sign in user (player/employee) with name, password (NO AUTH) --> RESPONSE: auth: "", istspieler
export const login = (name, password) => {
  return request({
    methode: 'login',
    param: {
      name,
      passwort: password,
    },
  })
}

//Sign out user (player/employee) with auth --> RESPONSE:
export const logout = (auth) => {
  return request({
    auth,
    methode: 'logout',
  })
}

//Change password of user (player/employee) with auth --> RESPONSE:
export const changePassword = (auth, password) => {
  return request({
    auth,
    methode: 'aendereKontodaten',
    param: {
      neuespasswort: password,
    },
  })
}

//Delete account of user (player/employee) with auth --> RESPONSE:
//also logs out (? TO TEST)
export const deleteAccount = (auth) => {
  return request({
    auth,
    methode: 'loescheKontodaten',
  })
}

//Submit tip made by player with date, auth -> RESPONSE:
//only if date in future & employee has set this date to be filled previously
export const submitTip = (auth, tip, date) => {
  return request({
    auth,
    methode: 'neuerTipp',
    param: {
      tipp: tip,
      datum: date,
    },
  })
}

//Show tips that player made with auth, OPTIONAL von, bis --> RESPONSE: statistik: []
export const getTips = (auth, from, until) => {
  return request({
    auth,
    methode: 'zeigeTipps',
    param: {
      von: from,
      bis: until,
    },
  })
}

//Show games open for players to enter/employees to see with auth --> RESPONSE: ziehungstage: []
export const getOpenGames = (auth) => {
  return request({
    auth,
    methode: 'zeigeAktuelleSpiele',
  })
}

//Show closed games with OPTIONAL date, with auth -> RESPONSE: statistik { [] datum, ziehung, auszahlungen {[] klasse, gewinner, gewinn}}
export const getClosedGames = (auth, from, until) => {
  return request({
    auth,
    methode: 'holeZiehungen',
    param: {
      von: from,
      bis: until,
    },
  })
}

//Open game by employee with auth, date --> RESPONSE:
//must be in future
export const openGame = (auth, date) => {
  return request({
    auth,
    methode: 'neueZiehung',
    param: {
      datum: date,
    },
  })
}

//Submit drawing made by employee with auth, draw, date --> RESPONSE:
//must be in future and for previously opened game
export const closeGame = (auth, date, draw) => {
  return request({
    auth,
    methode: 'beendeZiehung',
    param: {
      datum: date,
      ziehung: draw,
    },
  })
}
