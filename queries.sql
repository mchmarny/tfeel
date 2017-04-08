SELECT sum(r.score) as score_sum, t.by
FROM [mchmarny-dev:tfeel.tweets] t
JOIN [mchmarny-dev:tfeel.results] r on t.id = r.id
WHERE r.score < 0
GROUP BY t.by
ORDER BY 1

SELECT t.on, r.score, t.body
FROM [mchmarny-dev:tfeel.tweets] t
JOIN [mchmarny-dev:tfeel.results] r on t.id = r.id
WHERE t.by = 'itsdjbrad'
ORDER BY 1 desc
