/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng
 */

// Thirdparty imports
import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

// Component imports
import RatingRecord from '../components/rating_record';
import RecommendationIndex from '../components/recommender/recommendation_index';

// Store imports
import {
  allMoviesFetch,
  personalizedRecommendedMoviesFetch,
  recommendedMoviesFetch
} from '../store/movies/movie.action';

// Style imports
import './recommender.scss';


class Recommender extends React.Component {
  static propTypes = {
    session: PropTypes.object.isRequired,
    movies: PropTypes.object.isRequired,
    movieRatings: PropTypes.object.isRequired,
    movieYearRange: PropTypes.object.isRequired,
    moviePopularityPercentile: PropTypes.number.isRequired,
    dispatchAllMovieFetch: PropTypes.func.isRequired,
    dispatchPersonalizedRecommendedMoviesFetch: PropTypes.func.isRequired,
    dispatchRecommendedMoviesFetch: PropTypes.func.isRequired,
  };

  componentWillReceiveProps(nextProps) {
    let remainingRecommendedItems = nextProps.movies.recommended.size;

    this.props.movies.recommended.forEach((movieId) => {
      if (nextProps.movies.skipped.has(movieId) || nextProps.movies.rated.has(movieId)) {
        remainingRecommendedItems -= 1;
      }
    });

    const filterParamHasChanged = (
      this.props.movieYearRange.minYear !== nextProps.movieYearRange.minYear
      || this.props.movieYearRange.maxYear !== nextProps.movieYearRange.maxYear
      || this.props.moviePopularityPercentile !== nextProps.moviePopularityPercentile
    );

    let shouldFetchRecommendations = false;

    // If authenticated user receives stored movie ratings from backend
    if (Object.keys(this.props.movieRatings).length === 0 && Object.keys(nextProps.movieRatings).length >= 10) {
      shouldFetchRecommendations = true;
    }

    // Self-explanatory
    if (filterParamHasChanged) {
      shouldFetchRecommendations = true;
    }

    // If the frontend does not show any more recommendations to user
    if (remainingRecommendedItems === 0 && Object.keys(nextProps.movieRatings).length >= 10) {
      shouldFetchRecommendations = true;
    }

    if (shouldFetchRecommendations && this.props.session.currentUser === null) {
      this.props.dispatchRecommendedMoviesFetch(
        nextProps.movieYearRange,
        nextProps.moviePopularityPercentile,
        nextProps.movies.skipped,
        nextProps.movieRatings
      );
    } else if (shouldFetchRecommendations) {
      this.props.dispatchPersonalizedRecommendedMoviesFetch(
        this.props.session,
        nextProps.movieYearRange,
        nextProps.moviePopularityPercentile,
        nextProps.movies.skipped
      );
    }
  }

  componentDidMount() {
    if (Object.keys(this.props.movies.all).length === 0) {
      this.props.dispatchAllMovieFetch();
    }

    if (Object.keys(this.props.movieRatings).length >= 10) {
      if (this.props.session.currentUser !== null) {
        this.props.dispatchPersonalizedRecommendedMoviesFetch(
          this.props.session,
          this.props.movieYearRange,
          this.props.moviePopularityPercentile,
          this.props.movies.skipped
        );
      } else {
        this.props.dispatchRecommendedMoviesFetch(
          this.props.movieYearRange,
          this.props.moviePopularityPercentile,
          this.props.movies.skipped,
          this.props.movieRatings
        );
      }
    }
  }

  render() {
    if (Object.keys(this.props.movieRatings).length < 10) {
      return (
        <section className="recommender">
          <header>
            <h3>Recommendations are not ready</h3>
            <p>{"You haven't rated enough movies for the system to recommend you anything."}</p>
          </header>
        </section>
      );
    }

    return (
      <section className="recommender">
        <RecommendationIndex />
        <RatingRecord />
      </section>
    );
  }
}

const mapReduxStateToProps = (state) => ({
  session: state.session,
  movies: state.movies,
  movieRatings: state.movieRatings,
  movieYearRange: state.movieYearRange,
  moviePopularityPercentile: state.moviePopularityPercentile
});

const mapDispatchToProps = (dispatch) => ({
  dispatchAllMovieFetch: () => dispatch(allMoviesFetch()),
  dispatchPersonalizedRecommendedMoviesFetch: (session, yearRange, percentile, skipped) => {
    return dispatch(personalizedRecommendedMoviesFetch(session, yearRange, percentile, skipped));
  },
  dispatchRecommendedMoviesFetch: (yearRange, percentile, skipped, ratings) => {
    return dispatch(recommendedMoviesFetch(yearRange, percentile, skipped, ratings));
  }
});

export default connect(mapReduxStateToProps, mapDispatchToProps)(Recommender);
