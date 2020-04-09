module Main exposing (main)

import Browser exposing (Document, UrlRequest)
import Browser.Navigation as Nav
import Html exposing (..)
import Html.Attributes exposing (class, href, src, style)
import Json.Decode as Decode exposing (Value)
import Page.Home as Home
import Page.Journal as Journal
import Page.Links as Links
import Page.ListBooks as ListBooks
import Page.Login as Login
import Page.ViewBook as ViewBook
import Route exposing (Route)
import Session exposing (..)
import Url exposing (Url)


main : Program (Maybe String) Model Msg
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
    | HomePage
    | JournalPage Journal.Model
    | LinksPage Links.Model


type Msg
    = ListBooksPageMsg ListBooks.Msg
    | ViewBookPageMsg ViewBook.Msg
    | LoginPageMsg Login.Msg
    | JournalPageMsg Journal.Msg
    | LinksPageMsg Links.Msg
    | LinkClicked UrlRequest
    | UrlChanged Url
    | LoggedIn Msg


init : Maybe String -> Url -> Nav.Key -> ( Model, Cmd Msg )
init flags url navKey =
    let
        model =
            case flags of
                Nothing ->
                    { route = Route.parseUrl url
                    , page = NotFoundPage
                    , navKey = navKey
                    , session = Unauthenticated
                    }

                Just value ->
                    case Decode.decodeString Session.storageDecoder value of
                        Ok token ->
                            { route = Route.parseUrl url
                            , page = NotFoundPage
                            , navKey = navKey
                            , session = Authenticated token
                            }

                        _ ->
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

                ( Route.Home, Authenticated token ) ->
                    ( HomePage, Cmd.none )

                ( Route.Journal, Authenticated token ) ->
                    let
                        ( pageModel, pageCmds ) =
                            Journal.init token
                    in
                    ( JournalPage pageModel, Cmd.map JournalPageMsg pageCmds )

                ( Route.Links, Authenticated token ) ->
                    let
                        ( pageModel, pageCmds ) =
                            Links.init token
                    in
                    ( LinksPage pageModel, Cmd.map LinksPageMsg pageCmds )

                ( Route.Home, Unauthenticated ) ->
                    ( HomePage, Nav.pushUrl model.navKey "/login" )

                ( Route.Login, Unauthenticated ) ->
                    let
                        ( pageModel, pageCmds ) =
                            Login.init model.navKey
                    in
                    ( LoginPage pageModel, Cmd.map LoginPageMsg pageCmds )

                ( Route.Login, Authenticated token ) ->
                    ( HomePage, Nav.pushUrl model.navKey "/" )

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
    case model.session of
        Authenticated _ ->
            { title = "Alexandria"
            , body = [ viewHeader model, navbar, currentView model ]
            }

        Unauthenticated ->
            { title = "Alexandria"
            , body = [ viewHeader model, currentView model ]
            }


navbar : Html Msg
navbar =
    div []
        [ div [ class "tabs is-toggle is-centered transparent" ]
            [ ul []
                (List.map
                    (\area ->
                        li []
                            [ a [ href area.endpoint ]
                                [ span [ class "icon is-small" ]
                                    [ i [ class "fas", class area.icon ] []
                                    ]
                                , span [] [ text area.name ]
                                ]
                            ]
                    )
                    commonAreas
                )
            ]
        ]


type alias Area =
    { name : String
    , endpoint : String
    , icon : String
    }


commonAreas : List Area
commonAreas =
    [ { name = "Books"
      , endpoint = "/books"
      , icon = "fa-book"
      }
    , { name = "Papers"
      , endpoint = "/papers"
      , icon = "fa-paper-plane"
      }
    , { name = "Journal"
      , endpoint = "/journal"
      , icon = "fa-comment"
      }
    , { name = "Links"
      , endpoint = "/links"
      , icon = "fa-link"
      }
    , { name = "Talks"
      , endpoint = "/talks"
      , icon = "fa-video"
      }
    , { name = "Tags"
      , endpoint = "/tags"
      , icon = "fa-tags"
      }
    , { name = "Ideas"
      , endpoint = "/ideas"
      , icon = "fa-lightbulb"
      }
    , { name = "Map"
      , endpoint = "/map"
      , icon = "fa-project-diagram"
      }
    , { name = "Search"
      , endpoint = "/search"
      , icon = "fa-search"
      }
    ]


viewHeader : Model -> Html Msg
viewHeader model =
    div []
        [ nav [ class "navbar", class "is-dark" ]
            [ div [ class "navbar-brand" ]
                [ a [ href "/" ]
                    [ div [ class "navbar-item" ]
                        [ span [ style "color" "white" ] [ text "Alexandria", img [ src "/alexandria.png" ] [] ]
                        ]
                    ]
                ]
            ]
        ]


currentView : Model -> Html Msg
currentView model =
    case model.page of
        NotFoundPage ->
            notFoundView

        HomePage ->
            Home.view

        JournalPage pageModel ->
            Journal.view pageModel
                |> Html.map JournalPageMsg

        LinksPage pageModel ->
            Links.view pageModel
                |> Html.map LinksPageMsg

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

        ( JournalPageMsg subMsg, JournalPage pageModel ) ->
            let
                ( updatedPageModel, updatedCmd ) =
                    Journal.update subMsg pageModel
            in
            ( { model | page = JournalPage updatedPageModel }
            , Cmd.map JournalPageMsg updatedCmd
            )

        ( LinksPageMsg subMsg, LinksPage pageModel ) ->
            let
                ( updatedPageModel, updatedCmd ) =
                    Links.update subMsg pageModel
            in
            ( { model | page = LinksPage updatedPageModel }
            , Cmd.map LinksPageMsg updatedCmd
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
                            ( { model | session = Authenticated url.token }, Cmd.batch [ Session.storeCredWith url.token, Nav.pushUrl model.navKey "/" ] )

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
