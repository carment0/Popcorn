/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng
 */

// Thirdparty imports
import React from 'react';
import PropTypes from 'prop-types';
import Spinner from 'spin';
import Rating from 'react-rating';
import FlatButton from 'material-ui/FlatButton';

// Style
import './recommended_movie_item.scss';


class RecommendedMovieItem extends React.Component {
  static propTypes = {
    movieDetail: PropTypes.object,
    trailerList: PropTypes.array,
    session: PropTypes.object.isRequired,
    movieId: PropTypes.number.isRequired,
    imdbId: PropTypes.string.isRequired,
    playTrailer: PropTypes.func.isRequired,
    dispatchMovieSkip: PropTypes.func.isRequired,
    dispatchMovieRatingPost: PropTypes.func.isRequired,
    dispatchMovieRatingRecord: PropTypes.func.isRequired,
    dispatchMovieDetailFetch: PropTypes.func.isRequired,
    dispatchMovieTrailerFetch: PropTypes.func.isRequired
  }

  static defaultProps = {
    trailerList: [],
    movieDetail: {
      isDefaultProp: true,
      title: 'Loading...',
      year: '....',
      plot: 'Loading...'
    },
  };

  /**
   * Initiate fetching metadata of the current movie and spinning the loading bar so user knows that it is fetching
   * something.
   */
  fetchMovieDetail = () => {
    if (this.props.movieDetail.isDefaultProp) {
      this.spinner = new Spinner().spin();
      document.getElementById(this.props.movieId).appendChild(this.spinner.el);
      this.props.dispatchMovieDetailFetch(this.props.imdbId);
    }
  };

  /**
   * If user is authenticated, rating is saved to database otherwise the rating is cached in redux store and wait until
   * user sign ups and logs in.
   * @param {number} ratingValue
   */
  handleRatingSelect = (ratingValue) => {
    if (this.props.session.currentUser !== null) {
      this.props.dispatchMovieRatingPost(this.props.movieId, ratingValue);
    }

    this.props.dispatchMovieRatingRecord(this.props.movieId, ratingValue);
  };

  /**
   * Label a movie as skipped in Redux store.
   */
  handleMovieSkip = () => {
    this.props.dispatchMovieSkip(this.props.movieId);
  };

  /**
   * Play a movie trailer
   */
  handleTrailerPlay = () => {
    this.props.playTrailer(this.props.imdbId);
  };

  get poster() {
    if (this.props.movieDetail.poster) {
      return (
        <section className="poster">
          <img id={this.props.movieId} src={this.props.movieDetail.poster} alt="Movie Poster" />
        </section>
      );
    }

    return <section className="poster" />;
  }

  componentDidMount() {
    this.fetchMovieDetail();
    if (this.props.trailerList.length === 0) {
      this.props.dispatchMovieTrailerFetch(this.props.imdbId);
    }
  }

  componentWillReceiveProps(nextProps) {
    if (!nextProps.movieDetail.isDefaultProp && this.spinner) {
      this.spinner.stop();
    }
  }

  render() {
    const buttons = [];

    if (this.props.trailerList.length > 0) {
      buttons.push(<FlatButton key="play" label="Play Trailer" primary={true} onClick={this.handleTrailerPlay} />);
    }

    buttons.push(<FlatButton key="skip" label="Not Interested" secondary={true} onClick={this.handleMovieSkip} />);

    // The reason for nesting a content box inside recommended movie item is that the content box has a row flex
    // direction and the recommended movie item has a column flex direction. Also, adding padding from outer layer is
    // easier than to add margin.
    return (
      <article className="recommended movie-item" id={this.props.movieId}>
        <section className="content-container">
          <aside className="left">
            {this.poster}
          </aside>
          <aside className="right">
            <h3 className="title"><strong>{this.props.movieDetail.title}</strong>{this.props.movieDetail.year}</h3>
            <div className="plot">{this.props.movieDetail.plot}</div>
            <div className="rating-toolbar" >
              <div className="star-toolbar-container">
                <Rating fractions={2} onChange={this.handleRatingSelect} />
              </div>
              <div className="button-container">{buttons}</div>
            </div>
          </aside>
        </section>
      </article>
    );
  }
}

export default RecommendedMovieItem;
