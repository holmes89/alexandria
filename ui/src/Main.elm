module Main exposing (main)

import Browser exposing (Document, UrlRequest)
import Browser.Navigation as Nav
import Html exposing (..)
import Html.Attributes exposing (class, href, src, style)
import Page.ListBooks as ListBooks
import Page.Login as Login
import Page.ViewBook as ViewBook
import Route exposing (Route)
import Session exposing (..)
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
    , session : Session
    }


type Page
    = NotFoundPage
    | UnauthorizedPage
    | ListBooksPage ListBooks.Model
    | ViewBookPage ViewBook.Model
    | LoginPage Login.Model


type Msg
    = ListBooksPageMsg ListBooks.Msg
    | ViewBookPageMsg ViewBook.Msg
    | LoginPageMsg Login.Msg
    | LinkClicked UrlRequest
    | UrlChanged Url


init : () -> Url -> Nav.Key -> ( Model, Cmd Msg )
init flags url navKey =
    let
        model =
            { route = Route.parseUrl url
            , page = NotFoundPage
            , navKey = navKey
            , session = Unauthenticated
            }
    in
    initCurrentPage ( model, Cmd.none )


initCurrentPage : ( Model, Cmd Msg ) -> ( Model, Cmd Msg )
initCurrentPage ( model, existingCmds ) =
    let
        ( currentPage, mappedPageCmds ) =
            case ( model.route, model.session ) of
                ( Route.NotFound, _ ) ->
                    ( NotFoundPage, Cmd.none )

                ( Route.Login, Unauthenticated ) ->
                    let
                        ( pageModel, pageCmds ) =
                            Login.init model.navKey
                    in
                    ( LoginPage pageModel, Cmd.map LoginPageMsg pageCmds )

                ( Route.Login, Authenticated token ) ->
                    let
                        ( pageModel, pageCmds ) =
                            ListBooks.init token
                    in
                    ( ListBooksPage pageModel, Cmd.map ListBooksPageMsg pageCmds )

                ( Route.Books, Authenticated token ) ->
                    let
                        ( pageModel, pageCmds ) =
                            ListBooks.init token
                    in
                    ( ListBooksPage pageModel, Cmd.map ListBooksPageMsg pageCmds )

                ( Route.Book bookID, Authenticated token ) ->
                    let
                        ( pageModel, pageCmds ) =
                            ViewBook.init bookID model.navKey token
                    in
                    ( ViewBookPage pageModel, Cmd.map ViewBookPageMsg pageCmds )

                ( _, Unauthenticated ) ->
                    ( UnauthorizedPage, Cmd.none )
    in
    ( { model | page = currentPage }
    , Cmd.batch [ existingCmds, mappedPageCmds ]
    )


view : Model -> Document Msg
view model =
    { title = "Alexandria"
    , body = [ viewHeader model, currentView model ]
    }


viewHeader : Model -> Html Msg
viewHeader model =
    div []
        [ nav [ class "navbar", class "is-dark" ]
            [ div [ class "navbar-brand" ]
                [ div [ class "navbar-item" ]
                    [ span [] [ text "Alexandria", img [ src "/alexandria.png" ] [] ]
                    ]
                ]
            ]
        ]


currentView : Model -> Html Msg
currentView model =
    case model.page of
        NotFoundPage ->
            notFoundView

        UnauthorizedPage ->
            unauthorizedView

        LoginPage pageModel ->
            Login.view pageModel
                |> Html.map LoginPageMsg

        ListBooksPage pageModel ->
            ListBooks.view pageModel
                |> Html.map ListBooksPageMsg

        ViewBookPage pageModel ->
            ViewBook.view pageModel
                |> Html.map ViewBookPageMsg


notFoundView : Html msg
notFoundView =
    h3 [] [ text "Oops! The page you requested was not found!" ]


unauthorizedView : Html msg
unauthorizedView =
    h3 [] [ text "Forbidden" ]


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

        ( ViewBookPageMsg subMsg, ViewBookPage pageModel ) ->
            let
                ( updatedPageModel, updatedCmd ) =
                    ViewBook.update subMsg pageModel
            in
            ( { model | page = ViewBookPage updatedPageModel }
            , Cmd.map ViewBookPageMsg updatedCmd
            )

        ( LoginPageMsg subMsg, LoginPage pageModel ) ->
            let
                ( updatedPageModel, updatedCmd ) =
                    Login.update subMsg pageModel
            in
            case subMsg of
                Login.Login result ->
                    case result of
                        Ok url ->
                            ( { model | session = Authenticated url.token }, Nav.pushUrl model.navKey "/books" )

                        Err _ ->
                            ( { model | session = Unauthenticated }, Cmd.none )

                _ ->
                    ( { model | page = LoginPage updatedPageModel }
                    , Cmd.map LoginPageMsg updatedCmd
                    )

        ( LinkClicked urlRequest, _ ) ->
            case urlRequest of
                Browser.Internal url ->
                    ( model
                    , Nav.pushUrl model.navKey (Url.toString url)
                    )

                Browser.External url ->
                    ( model
                    , Nav.load url
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
