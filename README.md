# Popcorn
Wonder what movie to watch next with your friends and family? Use Popcorn to fetch movie recommendations!

## Local Development Setup
### Dependency
The default choice of Go is 1.9.3 for this project and we are using `dep` to perform package dependency management. For
frontend, we are using `npm` for JavaScript dependency management. Thus, you need the following:

* Go 1.9.3
* PostgreSQL 9.3
* Node 6.10
* npm 3.10
* dep latest

### Database
Step one, enter `psql` command line interface:
```
psql postgres
```

Check what the users you have in your local Postgres server:
```
\du
```

Create a user named `popcorn`
```
create user popcorn superuser createdb createrole;
```

Give your `popcorn` a password:
```
alter user popcorn with password 'popcorn';
```

Now you can go ahead and create a database with owner pointing to `popcorn`:
```
create database popcorn_development with owner=popcorn;
```

Quit `psql`:
```
\q
```

## Build Project
### Backend
Inside your popcorn directory, Run `go install` to build the binary for your server
```
go install
```

Run `go install` to build the binaries for other commands
```
go install ./cmd/...
```

To start the server, simply run
```
popcorn
```

To seed the database, simply run
```
seed
```

Check if your seeds are actually working
```
psql popcorn_development
```

Enter the following SQL command:
```
select title from movies;
```

### Frontend
Install all the required node modules
```
npm install
```

And then build your scripts and style sheets
```
npm run build:watch
```

## How Does Popcorn Work?
### Collaborative Filtering
Coming soon...

### Low Rank Matrix Factorization
Suppose that our system has **I** users and **J** movies. We assign **K** latent features to each user and movie in the
system. We can construct a matrix factorization as follows:

![factorization](./docs/factorization.png)

**X** represents the latent feature matrix of all users in our system. The greek letter **big theta** represents the
latent feature matrix for all movies in our system. The matrix product of user latent features and transpose of movie
latent features is the model predicted rating matrix.

<img src="https://latex.codecogs.com/png.latex?\dpi{150}&space;\bg_white&space;X\Theta^{T}&space;=&space;\hat{R}"
title="X\Theta^{T} = \hat{R}" />

Let **R** represents the the actual rating we received from the MovieLens dataset. For every missing value in **R**, we
will replace them with the average rating each movie received from the poll of users. Then we define the loss function as
follows:

<img src="https://latex.codecogs.com/png.latex?\dpi{150}&space;\bg_white&space;L_{X,&space;\Theta}&space;=&space;
\frac{1}{2}\Sigma_{i,j}&space;(X\Theta^{T}&space;-&space;R)^{2}&space;&plus;&space;\frac{\lambda}{2}\Sigma_{i,&space;k}
X^{2}&space;&plus;&space;\frac{\lambda}{2}\Sigma_{j,&space;k}\Theta^{2}" title="L_{X, \Theta} = \frac{1}{2}\Sigma_{i,j}
(X\Theta^{T} - R)^{2} + \frac{\lambda}{2}\Sigma_{i, k}X^{2} + \frac{\lambda}{2}\Sigma_{j, k}\Theta^{2}" />

Thus, figuring out the latent features for movies and users has become a constraint optimization problem.

### Partial Derivatives & Gradients
Let's find the gradient of **L** with respect to the output of our low-rank matrix model. The one-half term will get
cancel out by the square term when we take the derivatives.

<img src="https://latex.codecogs.com/png.latex?\dpi{150}&space;\bg_white&space;\frac{\partial&space;L}{\partial&space;
  \hat{R}}&space;=&space;\hat{R}&space;-&space;R" title="\frac{\partial L}{\partial \hat{R}} = \hat{R} - R" />

Now we proceed to seek the gradient of model output with respect to **X** and **big theta**.

<img src="https://latex.codecogs.com/png.latex?\dpi{150}&space;\bg_white&space;\frac{\partial&space;\hat{R}}
{\partial&space;X}&space;=&space;\Theta^{T}" title="\frac{\partial \hat{R}}{\partial X} = \Theta^{T}" />

<img src="https://latex.codecogs.com/png.latex?\dpi{150}&space;\bg_white&space;\frac{\partial&space;\hat{R}}
{\partial&space;\Theta}&space;=&space;X" title="\frac{\partial \hat{R}}{\partial \Theta} = X" />

Using chain rule, we can then derive the following results:

<img src="https://latex.codecogs.com/png.latex?\dpi{150}&space;\bg_white&space;\frac{\partial&space;L}{\partial&space;X}
&space;=&space;\frac{\partial&space;L}{\partial&space;\hat{R}}\frac{\partial&space;\hat{R}}{\partial&space;X}"
title="\frac{\partial L}{\partial X} = \frac{\partial L}{\partial \hat{R}}\frac{\partial \hat{R}}{\partial X}" />

<img src="https://latex.codecogs.com/png.latex?\dpi{150}&space;\bg_white&space;\frac{\partial&space;L}
{\partial&space;\Theta}&space;=&space;\frac{\partial&space;L}{\partial&space;\hat{R}}\frac{\partial&space;\hat{R}}
{\partial&space;\Theta}"
title="\frac{\partial L}{\partial \Theta} = \frac{\partial L}{\partial \hat{R}}\frac{\partial \hat{R}}{\partial \Theta}" />

In Python
```Python
"""Denote U as the user latent feature matrix and M as the movie latent feature matrix"""
model_pred = np.dot(U, M.T)

grad_pred = model_pred - R
grad_u = np.dot(grad_pred, M) + (reg * U)
grad_m = np.dot(grad_pred.T, U) + (reg * M)
```
