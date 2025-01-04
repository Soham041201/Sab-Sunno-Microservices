const MODEL = 'gpt-4o-mini-realtime-preview-2024-12-17';
const VOICE = 'verse';

const fetchEphimerialToken = async () => {
  const r = await fetch('https://api.openai.com/v1/realtime/sessions', {
    method: 'POST',
    headers: {
      Authorization: `Bearer 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`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      model: MODEL,
      voice: VOICE,
    }),
  });
  const data = await r.json();
  return data;
};

const sendLocalDescriptionToOpenAi  = async ({offer,token}) => {
    const r = await fetch('https://api.openai.com/v1/realtime?model='+MODEL, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token.client_secret.value}`,
          'Content-Type': 'application/sdp',
        },

        body: offer.sdp
        });
        const data = await r.text();
        return data;
}




module.exports = {fetchEphimerialToken,sendLocalDescriptionToOpenAi}
