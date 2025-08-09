<p>
    This project is an iteration of <a href="https://github.com/nem0z/wiki-pathfinder ">wiki-pathfinder<a>
    The goal is to produce an application capable of crawling, storing and indexing data to build a pathfinder throw wikipedia articles.
    To achive this, this project will use advanced tools like MQRabbit and Neo4j in the opposite of wiki-pathfinder where I will try to scale using internal tools. 
<p>

# Workflow

## Startup workflow
- Start the database
- Create DB tables / struct if doesn't exist
  
- Start Redis
- Load articles in Redis
  
- Start MQ Rabbit
- Create queues
- Load unprocessed articles from in the queue  

## App workflow

- Start the App
- Start the crawlers
- Start the handler for articles
- Start the handler for relations

### Crawler workflow

- Consume unprocessed articles queue
- Request the article's page
- Extract articles URL
- Push the result to queue article

In case of failure the crawler push the article to the back of the queue

### Article handler workflow

- Consume the queue articles
- Request the article id to Redis
- Create the aricle if doesn't exist
- Resolve all articles ids
- Push the result to the queue relations

### Article handler workflow

- Consume the queue relations
- Create all the links from parent_id to child_id
- If no error -> mark the parent article as processed