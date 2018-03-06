/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng
 */

// Thirdparty imports
import React from 'react';
import PropTypes from 'prop-types';
import { Carousel } from 'react-bootstrap';
import { CarouselItem } from 'react-bootstrap';

// Style imports
import './poster_slider.scss';


class PosterSlider extends React.Component {
  state = {
    posters: []
  };

  static propTypes = {
    movies: PropTypes.object.isRequired,
    movieDetails: PropTypes.object.isRequired
  };

  /**
   * @return {boolean}
   */
  get hasGottenAllPosters() {
    return this.state.posters.length === Object.keys(this.props.movies).length && this.state.posters.length !== 0;
  }

  /**
   * @return {Array}
   */
  get carouselItems() {
    const numberPerRow = 5;
    const carouselItems = [];
    for (let i = 0; i < this.state.posters.length; i += numberPerRow) {
      let j = i;
      const posterRow = [];
      while (j < i + numberPerRow && j < this.state.posters.length) {
        posterRow.push(<img alt="movie-poster" src={this.state.posters[j]} key={j} height={225} width={150} />);
        j += 1;
      }

      carouselItems.push(
        <CarouselItem key={i}>
          <section className="item-container">{posterRow}</section>
        </CarouselItem>
      );
    }

    return carouselItems;
  }

  componentWillReceiveProps(nextProps) {
    const posters = Object.keys(nextProps.movies)
      .filter((movieId) => {
        const movie = nextProps.movies[movieId];
        if (nextProps.movieDetails[movie.imdb_id]) {
          return true;
        }

        return false;
      })
      .map((movieId) => (nextProps.movieDetails[nextProps.movies[movieId].imdb_id].poster));

    this.setState({ posters });
  }


  render() {
    if (this.hasGottenAllPosters) {
      return (
        <section className="poster-slider">
          <Carousel interval={3000} style={{ backgroundColor: 'transparent' }}>
            {this.carouselItems}
          </Carousel>
        </section>
      );
    }

    return <section className="poster-slider" />;
  }
}

export default PosterSlider;
