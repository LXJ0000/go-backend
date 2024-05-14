wrk.method="GET"
wrk.headers["Authorization"]="Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoicm9vdCIsImlkIjoxODA5MTU2MDM5MzcxNjk0MDgsImV4cCI6MTcxNTY3MjA0MX0.FURJfVes4iBtC4D1wU0jUTwXtvlwP6HPjzJNPPIZCB8"
-- wrk -t4 -d5s -c50 -s ./script/wrk/post_info.lua http://localhost:8080/api/post?post_id=173255181931122688