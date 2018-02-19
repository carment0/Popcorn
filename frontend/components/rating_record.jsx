/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng
 */

// Thirdparty imports
import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

// Component imports
import RatedMovieItem from './rated_movie_item';

// Store imports
import { movieDetailFetch } from '../store/movies/detail.action';
import { recommendedMoviesFetch } from '../store/movies/movie.action';

// Style imports
import './rating_record.scss';


class RatingRecord extends React.Component {
  static propTypes = {
    session: PropTypes.object.isRequired,
    movies: PropTypes.object.isRequired,
    movieRatings: PropTypes.object.isRequired,
    movieDetails: PropTypes.object.isRequired,
    dispatchMovieDetailFetch: PropTypes.func.isRequired,
    dispatchRecommendedMoviesFetch: PropTypes.func.isRequired
  };

  get ratedMovies() {
    const ratedMovieIds = Object.keys(this.props.movieRatings);
    return ratedMovieIds.map((movieId) => {
      const movie = this.props.movies.all[movieId];
      return (
        <RatedMovieItem
          key={movie.id}
          movieId={movie.id}
          imdbId={movie.imdb_id}
          movieRating={this.props.movieRatings[movieId]}
          movieDetail={this.props.movieDetails[movie.imdb_id]}
          dispatchMovieDetailFetch={this.props.dispatchMovieDetailFetch} />
      );
    });
  }

  componentWillReceiveProps(nextProps) {
    if (Object.keys(nextProps.movieRatings).length < 10) {
      return;
    }

    if (Object.keys(nextProps.movieRatings).length !== Object.keys(this.props.movieRatings).length) {
      this.props.dispatchRecommendedMoviesFetch(this.props.session, nextProps.movieRatings);
    }
  }

  render() {
    return (
      <div className="rating-record-container">
        <div className="movies">
          {this.ratedMovies}
        </div>
      </div>
    );
  }
}

const mapReduxStateToProps = (state) => {
  return {
    session: state.session,
    movies: state.movies,
    movieRatings: state.movieRatings,
    movieDetails: state.movieDetails,
  };
};

const mapDispatchToProps = (dispatch) => {
  return {
    dispatchMovieDetailFetch: (imdbId) => dispatch(movieDetailFetch(imdbId)),
    dispatchRecommendedMoviesFetch: (session, ratings) => dispatch(recommendedMoviesFetch(session, ratings))
  };
};

export default connect(mapReduxStateToProps, mapDispatchToProps)(RatingRecord);
