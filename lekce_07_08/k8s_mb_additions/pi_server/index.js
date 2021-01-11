const express = require('express')
const app = express()
const port = 3000

app.get('/healthz', (req, res) => {
  res.send('Everything running, chief!')
})

app.get('/pi', (req, res) => {
  let i = 1n;
  let x = 3n * (10n ** 100020n);
  let pi = x;
  while (x > 0) {
    x = x * i / ((i + 1n) * 4n);
    pi += x / (i + 2n);
    i += 2n;
  }
  res.send((pi / (10n ** 20n)).toString())
})

app.listen(port, () => {
  console.log(`Pi server started`)
})


