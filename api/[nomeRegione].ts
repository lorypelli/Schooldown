import axios from 'axios';
export default async (request: import('@vercel/node').VercelRequest, response: import('@vercel/node').VercelResponse) => {
    response.setHeader('Access-Control-Allow-Origin', 'http://localhost:4321');
    if (request.method !== 'GET') {
        response.setHeader('Content-Type', 'text/plain');
        return response.status(405).send('Method not allowed');
    }
    else if (request.method === 'GET') {
        const res = await axios.get('https://schooldown.vercel.app/api/getData');
        if (res.status == 404) {
            response.setHeader('Content-Type', 'text/plain');
            return response.status(404).send('Not Found');
        }
        if (res.data[request.query.nomeRegione as string]) {
            return response.send(res.data[request.query.nomeRegione as string]);
        }
        else {
            const nomeRegione = Object.keys(res.data);
            const randomIndex = Math.floor(Math.random() * 21);
            response.setHeader('Content-Type', 'text/plain');
            response.status(400).send(nomeRegione[randomIndex]);
        }
    }
};