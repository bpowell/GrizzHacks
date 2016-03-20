# API For Database Data

## Endpoints:
1. /api/getall POST Requires: {ticker: ''}
2. /api/getrange POST Requires: {ticker: '', start: '123456789'} Optional: {end: '123456789'}
3. /api/getday POST Requires: {ticker: '', date: '2016-Mar-18'}
4. /api/gettickers GET
5. /api/getarticleids POST Requies: {ticker: ''}
6. /api/getarticle POST Requires: {id: ''}
7. /api/updatecount POST Requires: {id: '', count: ''}- id is for uniqueword id
8. /api/updateweights POST Requires: {id: '', weights: ''} - id is for uniqueword id
9. /api/adduniqueword POST Requires: {article_id: '', word: ''} Optional: {count: '', weights: ''}
10. /api/getwodsforarticle POST Requires: {article_id: ''}
11. /api/getinfoforword POST Requires: {word: ''}
