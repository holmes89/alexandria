module Main exposing (main)

import Browser exposing (Document, UrlRequest)
import Browser.Navigation as Nav
import Html exposing (..)
import Page.ListBooks as ListBooks
import Route exposing (Route)
import Url exposing (Url)


main : Program () Model Msg
main =
    Browser.application
        { init = init
        , view = view
        , update = update
        , subscriptions = \_ -> Sub.none
        , onUrlRequest = LinkClicked
        , onUrlChange = UrlChanged
        }


type alias Model =
    { route : Route
    , page : Page
    , navKey : Nav.Key
    }


type Page
    = NotFoundPage
    | ListBooksPage ListBooks.Model


type Msg
    = ListBooksPageMsg ListBooks.Msg
    | LinkClicked UrlRequest
    | UrlChanged Url


init : () -> Url -> Nav.Key -> ( Model, Cmd Msg )
init flags url navKey =
    let
        model =
            { route = Route.parseUrl url
            , page = NotFoundPage
            , navKey = navKey
            }
    in
    initCurrentPage ( model, Cmd.none )


initCurrentPage : ( Model, Cmd Msg ) -> ( Model, Cmd Msg )
initCurrentPage ( model, existingCmds ) =
    let
        ( currentPage, mappedPageCmds ) =
            case model.route of
                Route.NotFound ->
                    ( NotFoundPage, Cmd.none )

                Route.Books ->
                    let
                        ( pageModel, pageCmds ) =
                            ListBooks.init
                    in
                    ( ListBooksPage pageModel, Cmd.map ListBooksPageMsg pageCmds )
    in
    ( { model | page = currentPage }
    , Cmd.batch [ existingCmds, mappedPageCmds ]
    )


view : Model -> Document Msg
view model =
    { title = "Post App"
    , body = [ currentView model ]
    }


currentView : Model -> Html Msg
currentView model =
    case model.page of
        NotFoundPage ->
            notFoundView

        ListBooksPage pageModel ->
            ListBooks.view pageModel
                |> Html.map ListBooksPageMsg


notFoundView : Html msg
notFoundView =
    h3 [] [ text "Oops! The page you requested was not found!" ]


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case ( msg, model.page ) of
        ( ListBooksPageMsg subMsg, ListBooksPage pageModel ) ->
            let
                ( updatedPageModel, updatedCmd ) =
                    ListBooks.update subMsg pageModel
            in
            ( { model | page = ListBooksPage updatedPageModel }
            , Cmd.map ListBooksPageMsg updatedCmd
            )

        ( UrlChanged url, _ ) ->
            let
                newRoute =
                    Route.parseUrl url
            in
            ( { model | route = newRoute }, Cmd.none )
                |> initCurrentPage

        ( _, _ ) ->
            ( model, Cmd.none )
