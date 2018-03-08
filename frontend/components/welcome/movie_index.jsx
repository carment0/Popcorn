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
import { popularMoviesFetch } from '../../store/movies/movie.action';
import { movieDetailFetch } from '../../store/movies/detail.action';
import { movieRatingPost, movieRatingRecord } from '../../store/movies/rating.action';

// Style imports
import './movie_index.scss';


class MovieIndex extends React.Component {
  state = {
    displayMovies: {}
  };

  static propTypes = {
    errors: PropTypes.object.isRequired,
    session: PropTypes.object.isRequired,
    movies: PropTypes.object.isRequired,
    movieDetails: PropTypes.object.isRequired,
    dispatchMovieDetailFetch: PropTypes.func.isRequired,
    dispatchMovieRatingPost: PropTypes.func.isRequired,
    dispatchMovieRatingRecord: PropTypes.func.isRequired,
    dispatchPopularMoviesFetch: PropTypes.func.isRequired
  };

  /**
   * Grab more movies from either Redux store or ask the backend to populate Redux store.
   */
  handleClickMoreMovies = () => {
    if (this.props.movies.popular.size > 0 && Object.keys(this.props.movies.all).length > 0) {
      this.shuffleMoviesAndSetDisplay(this.props.movies.all, this.props.movies.popular);
    } else {
      this.props.dispatchPopularMoviesFetch();
    }
  };

  /**
   * Either set data from this.props or nextProps to state, it shuffles the set of popular movies and render them to display
   * @param {Object} movies A dictionary of all movies
   * @param {Set} popularSet A set of popular movie ID
   */
  shuffleMoviesAndSetDisplay = (movies, popularSet) => {
    if (popularSet.size > 0 && Object.keys(movies).length > 0) {
      const displayMovies = {};
      _.shuffle(Array.from(popularSet)).slice(0, 20).forEach((movieId) => {
        const movie = movies[movieId];
        displayMovies[movieId] = movie;
      });

      this.setState({ displayMovies });
    }
  };

  /**
   * Returns an instruction section
   * @return {React.Element}
   */
  get instruction() {
    if (Object.keys(this.state.displayMovies).length === 0) {
      return <div className="instruction" />;
    }

    const instruction = `These are some of the most popular American films. We think it is very likely that you have
    seen at least some of them.  If you have seen them, whether you like or dislike them, let us know and give them
    ratings! It will help our backend machine learning algorithm to learn your taste and preference`;
    const ratingCount = this.props.movies.rated.size;
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

  /**
   * Returns a list of movie items that are in the popular list and not yet rated.
   * @return {Array}
   */
  get popularMovieItems() {
    // NOTE: movieId is ALWAYS integer except when you do Object.keys()
    const unratedMovieIds = Object.keys(this.state.displayMovies).filter((movieId) => {
      const movie = this.state.displayMovies[movieId];
      const isRated = this.props.movies.rated.has(movie.id);
      const isDetailFetchSuccess = this.props.errors.movieDetails[movie.imdb_id] === undefined;
      return !isRated && isDetailFetchSuccess;
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
          dispatchMovieRatingPost={this.props.dispatchMovieRatingPost}
          dispatchMovieRatingRecord={this.props.dispatchMovieRatingRecord}
          dispatchMovieDetailFetch={this.props.dispatchMovieDetailFetch} />
      );
    });
  }

  /**
   * Returns a header section
   * @return {React.Element}
   */
  get header() {
    if (Object.keys(this.state.displayMovies).length > 0) {
      return (
        <section className="header">
          <h1>Popular Movies</h1>
        </section>
      );
    }

    return (
      <section className="header">
        <h1>Popular Movies</h1>
        <p>{"If you don't see any movie showing below, just click the button to load more movies."}</p>
      </section>
    );
  }

  componentDidMount() {
    if (this.props.movies.popular.size === 0) {
      this.props.dispatchPopularMoviesFetch();
    } else {
      this.shuffleMoviesAndSetDisplay(this.props.movies.all, this.props.movies.popular);
    }
  }

  componentWillReceiveProps(nextProps) {
    const popularMoviesPropHasChanged = this.props.movies.popular.size !== nextProps.movies.popular.size;

    const moviesPropHasChanged = (
      Object.keys(this.props.movies.all).length !== Object.keys(nextProps.movies.all).length
    );

    if (popularMoviesPropHasChanged || moviesPropHasChanged) {
      this.shuffleMoviesAndSetDisplay(nextProps.movies.all, nextProps.movies.popular);
    }
  }

  render() {
    const progressPercentage = (100 * this.props.movies.rated.size) / 10;
    if (
      this.props.movies.popular.size === 0
      || Object.keys(this.props.movies.all).length === 0
      || Object.keys(this.state.displayMovies).length === 0
    ) {
      return (
        <article className="movie-index">
          <section className="header">
            <h1>Popular Movies</h1>
          </section>
          <section className="fetching-progress">
            <h5>Please wait while movies are fetching; Heroku takes time to start up its server.</h5>
            <LinearProgress mode="indeterminate" />
          </section>
        </article>
      );
    }

    return (
      <article className="movie-index">
        {this.header}
        <PosterSlider
          movies={this.state.displayMovies}
          movieDetails={this.props.movieDetails} />
        {this.instruction}
        <LinearProgress mode="determinate" value={progressPercentage} />
        <div className="movies">
          { this.popularMovieItems }
        </div>
        <section className="footer">
          <Button bsSize="xsmall" className="react-buttons" onClick={this.handleClickMoreMovies} bsStyle="primary">
            Load more movies
          </Button>
        </section>
      </article>
    );
  }
}

const mapReduxStateToProps = (state) => {
  return {
    errors: state.errors,
    movies: state.movies,
    movieDetails: state.movieDetails,
    session: state.session
  };
};

const mapDispatchToProps = (dispatch) => {
  return {
    dispatchPopularMoviesFetch: () => dispatch(popularMoviesFetch()),
    dispatchMovieDetailFetch: (imdbId) => dispatch(movieDetailFetch(imdbId)),
    dispatchMovieRatingRecord: (movieId, rating) => dispatch(movieRatingRecord(movieId, rating)),
    dispatchMovieRatingPost: (movieId, userId, rating) => dispatch(movieRatingPost(movieId, userId, rating))
  };
};

export default connect(mapReduxStateToProps, mapDispatchToProps)(MovieIndex);
