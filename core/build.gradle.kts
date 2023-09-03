import org.apache.tools.ant.filters.Native2AsciiFilter
import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    kotlin("jvm") version "1.9.10"
    kotlin("kapt") version "1.9.10"

    application
}

group = "com.wittano.komputer"

repositories {
    mavenCentral()
}

val jacksonVersion = "2.15.2"
val daggerVersion = "2.48"

dependencies {
    // Discord4j
    implementation("com.discord4j:discord4j-core:3.2.5")
    implementation("io.projectreactor.kotlin:reactor-kotlin-extensions:1.2.2")

    // Logger
    implementation("ch.qos.logback:logback-classic:1.4.11")
    implementation("org.codehaus.janino:janino:3.1.10")

    // Dagger
    implementation("com.google.dagger:dagger:$daggerVersion")
    kapt("com.google.dagger:dagger-compiler:$daggerVersion")

    // MongoDB
    implementation("org.mongodb:mongodb-driver-reactivestreams:4.10.0")

    // Utils
    implementation("io.github.cdimascio:dotenv-kotlin:6.4.1")
    implementation("com.squareup.okhttp3:okhttp:4.11.0")

    // Jackson object mapper
    implementation("com.fasterxml.jackson.core:jackson-core:$jacksonVersion")
    implementation("com.fasterxml.jackson.module:jackson-module-kotlin:$jacksonVersion")

    testImplementation(platform("org.junit:junit-bom:5.9.1"))
    testImplementation("org.junit.jupiter:junit-jupiter")
}

kapt {
    arguments {
        arg("project", "${project.group}/${project.name}")
    }
}

tasks.withType<Jar> {
    manifest {
        attributes["Main-Class"] = "com.wittano.komputer.MainKt"
    }

    duplicatesStrategy = DuplicatesStrategy.EXCLUDE

    from(sourceSets.main.get().output)

    dependsOn(configurations.runtimeClasspath)
    from({
        configurations.runtimeClasspath.get().filter { it.name.endsWith("jar") }.map { zipTree(it) }
    })
}

val native2Ascii = Native2AsciiFilter()

tasks.withType<ProcessResources>().configureEach {
    filesMatching("**/i18n/*.properties") {
        filter {
            native2Ascii.filter(it)
        }
    }
}

tasks.test {
    useJUnitPlatform()
}

tasks.withType<KotlinCompile> {
    kotlinOptions.jvmTarget = "17"
}