/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng
 */

// Thirdparty imports
import _ from 'lodash';
import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Button } from 'react-bootstrap';
import LinearProgress from 'material-ui/LinearProgress';

// Component imports
import MovieItem from './movie_item';
import PosterSlider from './poster_slider';

// Store imports
import { popularMoviesFetch, movieSkip } from '../../store/movies/movie.action';
import { movieDetailFetch, movieTrailerFetch } from '../../store/movies/detail.action';
import { movieRatingPost, movieRatingRecord } from '../../store/movies/rating.action';

// Style imports
import './movie_index.scss';


class MovieIndex extends React.Component {
  state = {
    displayMovies: {}
  };

  static propTypes = {
    session: PropTypes.object.isRequired,
    movies: PropTypes.object.isRequired,
    movieDetails: PropTypes.object.isRequired,
    movieRatings: PropTypes.object.isRequired,
    dispatchMovieSkip: PropTypes.func.isRequired,
    dispatchMovieDetailFetch: PropTypes.func.isRequired,
    dispatchMovieRatingPost: PropTypes.func.isRequired,
    dispatchMovieRatingRecord: PropTypes.func.isRequired,
    dispatchPopularMoviesFetch: PropTypes.func.isRequired
  };

  handleButtonClickMoreMovies = () => {
    if (this.props.movies.popular.size > 0 && Object.keys(this.props.movies.all).length > 0) {
      this.shuffleMoviesAndSetDisplay(this.props.movies.all, this.props.movies.popular);
    } else {
      this.props.dispatchPopularMoviesFetch();
    }
  };

  shuffleMoviesAndSetDisplay = (movies, popularSet) => {
    if (popularSet.size > 0 && Object.keys(movies).length > 0) {
      const displayMovies = {};
      _.shuffle(Array.from(popularSet)).slice(0, 12).forEach((movieId) => {
        const movie = movies[movieId];
        displayMovies[movieId] = movie;
      });

      this.setState({ displayMovies });
    }
  };

  get instruction() {
    const instruction = `These are some of the most popular American films. We think it is very likely that you have
    seen at least some of them.  If you have seen them, whether you like or dislike them, let us know and give them
    ratings! It will help our backend machine learning algorithm to learn your taste and preference`;
    const ratingCount = Object.keys(this.props.movieRatings).length;
    if (ratingCount === 0) {
      return (
        <div className="instruction">
          <p>{instruction}</p>
          <h4>As soon as 10 movies are rated, the recommender system will get to work!</h4>
        </div>
      );
    } else if (ratingCount >= 10) {
      return (
        <div className="instruction">
          <p>{instruction}</p>
          <h4>Recommendations are ready!</h4>
        </div>
      );
    }

    return (
      <div className="instruction">
        <p>{instruction}</p>
        <h4>Rate <strong>{10 - ratingCount}</strong> more movies</h4>
      </div>
    );
  }

  get popularMovieItems() {
    // Filter away the rated items.
    const unratedMovieIds = Object.keys(this.state.displayMovies).filter((movieId) => {
      return !this.props.movieRatings[movieId];
    });

    return unratedMovieIds.sort().map((movieId) => {
      const movie = this.state.displayMovies[movieId];
      return (
        <MovieItem
          key={movie.id}
          movieDetail={this.props.movieDetails[movie.imdb_id]}
          movieId={movie.id}
          imdbId={movie.imdb_id}
          session={this.props.session}
          dispatchMovieSkip={this.props.dispatchMovieSkip}
          dispatchMovieRatingPost={this.props.dispatchMovieRatingPost}
          dispatchMovieRatingRecord={this.props.dispatchMovieRatingRecord}
          dispatchMovieDetailFetch={this.props.dispatchMovieDetailFetch} />
      );
    });
  }

  componentDidMount() {
    if (Object.keys(this.props.movies.popular).length === 0) {
      this.props.dispatchPopularMoviesFetch();
    }
  }

  componentWillReceiveProps(nextProps) {
    const popularMoviesPropHasChanged = this.props.movies.popular.size !== nextProps.movies.popular.size;

    const moviesPropHasChanged = (
      Object.keys(this.props.movies.all).length !== Object.keys(nextProps.movies.all).length
    );

    if (popularMoviesPropHasChanged || moviesPropHasChanged) {
      console.log('Prop has changed, now shuffle movie and display!');
      this.shuffleMoviesAndSetDisplay(nextProps.movies.all, nextProps.movies.popular);
    }
  }

  render() {
    const progressPercentage = (100 * Object.keys(this.props.movieRatings).length) / 10;
    if (this.props.movies.popular.size === 0 || Object.keys(this.props.movies.all) === 0) {
      return (
        <article className="movie-index">
          <section className="header">
            <h1>Popular Movies</h1>
          </section>
          <section className="fetching-progress">
            <h5>Please wait while movies are fetching</h5>
            <LinearProgress mode="indeterminate" />
          </section>
        </article>
      );
    }

    return (
      <article className="movie-index">
        <header>
          <h1>Popular Movies</h1>
        </header>
        <PosterSlider
          movies={this.state.displayMovies}
          movieDetails={this.props.movieDetails} />
        {this.instruction}
        <LinearProgress mode="determinate" value={progressPercentage} />
        <div className="movies">
          { this.popularMovieItems }
        </div>
        <section className="footer">
          <Button
            bsSize="xsmall"
            className="react-buttons"
            onClick={this.handleButtonClickMoreMovies}
            bsStyle="primary">
            Load more movies
          </Button>
        </section>
      </article>
    );
  }
}

const mapReduxStateToProps = (state) => {
  return {
    movies: state.movies,
    movieDetails: state.movieDetails,
    movieRatings: state.movieRatings,
    session: state.session
  };
};

const mapDispatchToProps = (dispatch) => {
  return {
    dispatchPopularMoviesFetch: () => dispatch(popularMoviesFetch()),
    dispatchMovieSkip: (movieId) => dispatch(movieSkip(movieId)),
    dispatchMovieDetailFetch: (imdbId) => dispatch(movieDetailFetch(imdbId)),
    dispatchMovieTrailerFetch: (imdbId) => dispatch(movieTrailerFetch(imdbId)),
    dispatchMovieRatingRecord: (movieId, rating) => dispatch(movieRatingRecord(movieId, rating)),
    dispatchMovieRatingPost: (movieId, rating) => dispatch(movieRatingPost(movieId, rating))
  };
};

export default connect(mapReduxStateToProps, mapDispatchToProps)(MovieIndex);
